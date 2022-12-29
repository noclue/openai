package main

import (
	"context"
	"fmt"
	"os"

	"github.com/noclue/openai"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "openai",
		Short: "OpenAI CLI",
		Long:  `OpenAI CLI provides command line tools for interacting with the OpenAI API. To authorize access set the OPENAI_API_KEY environment variable to your OpenAI API key.`,
	}

	var n int
	var size string
	var responseFormat string
	var user string
	var mask string
	var inputFile string
	var temperature float64
	var topP float64
	var outputFile string

	addImageFlags := func(cmd *cobra.Command) {
		cmd.Flags().IntVarP(&n, "num-images", "n", 1, "number of images (optional, default: 1)")
		cmd.Flags().StringVarP(&size, "size", "s", "medium", "size (small, medium, large) (optional, default: medium)")
		cmd.Flags().StringVarP(&responseFormat, "response-format", "r", "url", "response format (url, b64_json) (optional, default: url)")
		cmd.Flags().StringVarP(&user, "user", "u", "", "user (optional, default: none)")
	}

	var imageCmd = &cobra.Command{
		Use:   "image",
		Short: "Given a prompt and/or an input image, the model will generate a new image.",
		Long:  `Given a prompt and/or an input image, the model will generate a new image.`,
	}
	rootCmd.AddCommand(imageCmd)

	var createImageCmd = &cobra.Command{
		Use:   "create [prompt]",
		Short: "Create an image given a prompt",
		Long:  `Create an image given a prompt. The prompt is a text description of the desired image(s). The maximum length of the prmpt is 1000 characters.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			prompt := args[0]
			createImage(prompt, n, getSize(size), getResponseFormat(responseFormat), user)
		},
	}
	addImageFlags(createImageCmd)
	imageCmd.AddCommand(createImageCmd)

	var createImageVariationsCmd = &cobra.Command{
		Use:   "variations [image file]",
		Short: "Create image variations from the provided image",
		Long:  `Create image variations from the provided image. The image to use as the basis for the variation(s) must be a valid PNG file, less than 4MB, and square.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			imageFile := args[0]
			imageVariations(imageFile, n, getSize(size), getResponseFormat(responseFormat), user)
		},
	}
	addImageFlags(createImageVariationsCmd)
	imageCmd.AddCommand(createImageVariationsCmd)

	var CreateImageEditsCmd = &cobra.Command{
		Use:   "edits [image file] [prompt]",
		Short: "Create image edits from the provided image and prompt",
		Long:  `Create image edits from the provided image and prompt. The image to use as the basis for the edit(s) must be a valid PNG file, less than 4MB, and square. The prompt is a text description of the desired edit(s). The maximum length of the prmpt is 1000 characters.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			imageFile := args[0]
			prompt := args[1]
			imageEdits(imageFile, prompt, mask, n, getSize(size), getResponseFormat(responseFormat), user)
		},
	}
	addImageFlags(CreateImageEditsCmd)
	CreateImageEditsCmd.Flags().StringVarP(&mask, "mask", "m", "", "An additional image whose fully transparent areas (e.g. where alpha is zero) indicate where image should be edited. Must be a valid PNG file, less than 4MB, and have the same dimensions as image. (optional, default: none)")
	imageCmd.AddCommand(CreateImageEditsCmd)

	// Edit command
	var editCmd = &cobra.Command{
		Use:   "edit [model] [instruction file]",
		Short: "Edit a text using the specified model and instruction",
		Long:  `Edit a text using the specified model and instruction. The model is the name of the model to use. The instruction is a text description of the desired edit(s). The maximum length of the instruction is 1000 characters.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			model := args[0]
			instructionFile := args[1]
			edit(model, instructionFile, inputFile, n, temperature, topP, outputFile)
		},
	}
	editCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "input file (optional, default: none)")
	editCmd.Flags().IntVarP(&n, "num-edits", "n", 1, "number of edits (optional, default: 1)")
	editCmd.Flags().Float64VarP(&temperature, "temperature", "t", 1.0, "temperature (optional, default: 1.0)")
	editCmd.Flags().Float64VarP(&topP, "top-p", "p", 1.0, "top p (optional, default: 1.0)")
	editCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "output file (optional, default: none)")
	rootCmd.AddCommand(editCmd)

	rootCmd.Execute()

}

func getResponseFormat(responseFormat string) openai.ResponseFormat {
	switch responseFormat {
	case "url":
		return openai.Url
	case "b64_json":
		return openai.B64_json
	default:
		fmt.Printf("Invalid response format: %s\n", responseFormat)
		os.Exit(1)
		return openai.Url
	}
}

func getSize(size string) openai.ImageSize {
	switch size {
	case "small":
		return openai.SmallImage
	case "medium":
		return openai.MediumImage
	case "large":
		return openai.LargeImage
	default:
		fmt.Printf("Invalid size: %s\n", size)
		os.Exit(1)
		return openai.MediumImage
	}
}

// createImage creates an image
func createImage(prompt string, numImages int, size openai.ImageSize, responseFormat openai.ResponseFormat, user string) {
	if len(prompt) < 5 {
		fmt.Println("Prompt must be at least 5 characters long")
		os.Exit(1)
	}
	client := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	res, err := client.CreateImage(context.Background(), openai.CreateImageReq{
		Prompt: prompt,
		CommonImageReq: openai.CommonImageReq{
			N:              &numImages,
			Size:           size,
			ResponseFormat: responseFormat,
			User:           user,
		},
	})

	if err != nil {
		fmt.Printf("Encountered Error: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %+v", res)
}

// imageVariations creates variations of an image
func imageVariations(imageFile string, numImages int, size openai.ImageSize, responseFormat openai.ResponseFormat, user string) {
	if _, err := os.Stat(imageFile); os.IsNotExist(err) {
		fmt.Println("Image file does not exist: ", imageFile)
		os.Exit(1)
	}
	client := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	res, err := client.CreateImageVariations(context.Background(), openai.CreateImageVariationsReq{
		Image: imageFile,
		CommonImageReq: openai.CommonImageReq{
			N:              &numImages,
			Size:           size,
			ResponseFormat: responseFormat,
			User:           user,
		},
	})

	if err != nil {
		fmt.Printf("Encountered Error: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %+v", res)
}

func imageEdits(imageFile string, prompt string, mask string, numImages int, size openai.ImageSize, responseFormat openai.ResponseFormat, user string) {
	if _, err := os.Stat(imageFile); os.IsNotExist(err) {
		fmt.Println("Image file does not exist: ", imageFile)
		os.Exit(1)
	}
	if mask != "" {
		if _, err := os.Stat(mask); os.IsNotExist(err) {
			fmt.Println("Mask file does not exist: ", mask)
			os.Exit(1)
		}
	}
	if len(prompt) < 5 {
		fmt.Println("Prompt must be at least 5 characters long")
		os.Exit(1)
	}
	client := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	res, err := client.CreateImageEdits(context.Background(), openai.CreateImageEditsReq{
		Image:  imageFile,
		Prompt: prompt,
		Mask:   mask,
		CommonImageReq: openai.CommonImageReq{
			N:              &numImages,
			Size:           size,
			ResponseFormat: responseFormat,
			User:           user,
		},
	})

	if err != nil {
		fmt.Printf("Encountered Error: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %+v", res)
}

func edit(model string, instructionFile string, inputFile string, n int, temperature float64, topP float64, outputFile string) {
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
	instruction := string(instructionBytes)

	var input string
	if inputFile != "" {
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
		fmt.Printf("Response: %+v", res)
	}
}
