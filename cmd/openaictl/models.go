package openaictl

import (
	"context"
	"fmt"
	"os"

	"github.com/noclue/openai"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func addModelsCmd(rootCmd *cobra.Command) {
	var modelsCmd = &cobra.Command{
		Use:   "models",
		Short: "List models",
		Long:  `List OpenAI supported models as yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			models()
		},
	}
	rootCmd.AddCommand(modelsCmd)
}

func models() {
	client := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	resp, err := client.Models(context.Background())
	if err != nil {
		fmt.Printf("Error listing models: %s", err)
		os.Exit(1)
	}
	yaml, err := yaml.Marshal(resp)
	if err != nil {
		fmt.Printf("Error marshalling response to yaml: %s", err)
		os.Exit(1)
	}
	fmt.Println(string(yaml))
}
