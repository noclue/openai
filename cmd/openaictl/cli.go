package openaictl

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Flags
var n int
var size string
var responseFormat string
var user string
var mask string
var inputFile string
var input string
var temperature float64
var topP float64
var outputFile string
var model string
var instructionFile string
var instruction string

func Run() {
	var rootCmd = &cobra.Command{
		Use:   "openai",
		Short: "OpenAI CLI",
		Long:  `OpenAI CLI is a command line tool for interacting with the OpenAI API. To authorize access set the OPENAI_API_KEY environment variable to your OpenAI API key.`,
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(imageCmd())

	rootCmd.AddCommand(editCmd())

	rootCmd.AddCommand(modelsCmd())

	rootCmd.AddCommand(moderationsCmd())

	rootCmd.Execute()

}

// printResponse prints the response as yaml
func printResponse(res any) {
	y, err := yaml.Marshal(res)
	if err != nil {
		fmt.Printf("Error marshalling response to yaml: %+v", err)
		os.Exit(1)
	}
	fmt.Println(string(y))
}
