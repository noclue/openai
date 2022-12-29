package openaictl

import (
	"github.com/spf13/cobra"
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
		Long:  `OpenAI CLI provides command line tools for interacting with the OpenAI API. To authorize access set the OPENAI_API_KEY environment variable to your OpenAI API key.`,
	}
	addImageCmd(rootCmd)

	addEditCmd(rootCmd)

	addModelsCmd(rootCmd)

	rootCmd.Execute()

}
