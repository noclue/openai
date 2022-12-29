package openaictl

import (
	"context"
	"fmt"
	"os"

	"github.com/noclue/openai"
	"github.com/spf13/cobra"
)

func addEditCmd(rootCmd *cobra.Command) {
	var editCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit input text using the specified model and instruction",
		Long:  `Edit input text using the specified model and instruction. The model is the name of the model to use. The instruction is a text description of the desired edit(s). The maximum length of the instruction is 1000 characters. The input is the text to be edited. The number of edits is the number of edits to make. The temperature is a number between 0 and 1 that controls the randomness of the edits. The top p is a number between 0 and 1 that controls the diversity of the edits. The output file is where the edited text will be written.`,
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			edit(model, instructionFile, instruction, inputFile, input, n, temperature, topP, outputFile)
		},
	}
	editCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file (optional, default: none)")
	editCmd.Flags().StringVarP(&input, "input", "a", "", "input (optional, default: none)")
	editCmd.Flags().IntVarP(&n, "num-edits", "n", 1, "number of edits (optional, default: 1)")
	editCmd.Flags().Float64VarP(&temperature, "temperature", "t", 1.0, "temperature (optional, default: 1.0)")
	editCmd.Flags().Float64VarP(&topP, "top-p", "p", 1.0, "top p (optional, default: 1.0)")
	editCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "output file (optional, default: none)")
	editCmd.Flags().StringVarP(&model, "model", "m", "text-davinci-edit-001", "model (optional, default: text-davinci-edit-001)")
	editCmd.Flags().StringVarP(&instructionFile, "instruction-file", "f", "", "instruction file (required)")
	editCmd.Flags().StringVarP(&instruction, "instruction", "s", "", "instruction (optional, default: none)")
	rootCmd.AddCommand(editCmd)
}

func edit(model, instructionFile, instruction, inputFile, input string, n int, temperature float64, topP float64, outputFile string) {
	if instructionFile == "" && instruction == "" {
		fmt.Println("Instruction or instruction file is required")
		os.Exit(1)
	} else if instructionFile != "" && instruction != "" {
		fmt.Println("Instruction and instruction file are mutually exclusive")
		os.Exit(1)
	} else if instructionFile != "" {
		if _, err := os.Stat(instructionFile); os.IsNotExist(err) {
			fmt.Println("Instruction file does not exist: ", instructionFile)
			os.Exit(1)
		}
		// Read instruction from instructionFile
		instructionBytes, err := os.ReadFile(instructionFile)
		if err != nil {
			fmt.Printf("Error reading instruction file: %s", err)
			os.Exit(1)
		}
		instruction = string(instructionBytes)
	}

	if inputFile != "" && input != "" {
		fmt.Println("Input and input file are mutually exclusive")
		os.Exit(1)
	} else if inputFile != "" {
		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Println("Input file does not exist: ", inputFile)
			os.Exit(1)
		}
		inputBytes, err := os.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("Error reading input file: %s", err)
			os.Exit(1)
		}
		input = string(inputBytes)
	}

	client := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	params := openai.EditRequest{
		Model:       model,
		Instruction: instruction,
		Input:       input,
	}
	if n != 1 {
		params.N = &n
	}
	if temperature != 1.0 {
		params.Temperature = &temperature
	}
	if topP != 1.0 {
		params.TopP = &topP
	}
	res, err := client.Edit(context.Background(), params)

	if err != nil {
		fmt.Printf("Encountered Error: %+v", err)
		os.Exit(1)
	}

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(res.Choices[0].Text), 0644)
		if err != nil {
			fmt.Printf("Error writing output file: %s", err)
			os.Exit(1)
		}
	} else {
		printResponse(res)
	}
}
