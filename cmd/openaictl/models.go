package openaictl

import (
	"context"
	"fmt"
	"os"

	"github.com/noclue/openai"
	"github.com/spf13/cobra"
)

func modelsCmd() *cobra.Command {
	var modelsCmd = &cobra.Command{
		Use:   "models",
		Short: "List models",
		Long:  `List OpenAI supported models as yaml`,
		Run: func(cmd *cobra.Command, args []string) {
			models()
		},
	}
	return modelsCmd
}

func models() {
	client := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	res, err := client.Models(context.Background())
	if err != nil {
		fmt.Printf("Error listing models: %s", err)
		os.Exit(1)
	}
	printResponse(res)
}
