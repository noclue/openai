package openaictl

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/noclue/openai"
	"github.com/spf13/cobra"
)

// addModerationsCmd adds the moderations command to the root command.
func addModerationsCmd(rootCmd *cobra.Command) {
	var moderationsCmd = &cobra.Command{
		Use:   "moderation [flags]",
		Short: "Given a input text, outputs if the model classifies it as violating OpenAI's content policy.",
		Long:  `Given a input text, outputs if the model classifies it as violating OpenAI's content policy.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			moderations()
		},
	}
	moderationsCmd.Flags().StringVarP(&input, "input", "i", "", "The text to be moderated")
	moderationsCmd.Flags().StringVarP(&inputFile, "input-file", "f", "", "The file containing the text to be moderated")
	moderationsCmd.Flags().StringVarP(&model, "model", "m", "", "The model to use. Defaults to text-moderation-latest")
	rootCmd.AddCommand(moderationsCmd)
}

// moderations runs the moderations command.
func moderations() {
	if inputFile != "" && input != "" {
		fmt.Println("Input and input file are mutually exclusive")
		os.Exit(1)
	} else if inputFile != "" {
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Printf("Input file %s does not exist", inputFile)
			os.Exit(1)
		}
		inputBytes, err := ioutil.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading input file: %+v", err)
			os.Exit(1)
		}
		input = string(inputBytes)
	} else if input == "" {
		fmt.Println("Input or input file is required")
		os.Exit(1)
	}
	params := openai.ModerationRequest{
		Input: []string{input},
	}
	if model != "" {
		params.Model = model
	}
	c := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	res, err := c.Moderation(context.Background(), params)
	if err != nil {
		fmt.Printf("Error calling moderation: %+v", err)
		os.Exit(1)
	}
	printResponse(res)
}
