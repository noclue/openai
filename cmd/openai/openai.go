package main

import (
	"fmt"
	"os"

	"github.com/noclue/openai/openai"
	"github.com/spf13/cobra"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "openai",
		Short: "OpenAI CLI",
		Long:  `OpenAI CLI provides command line tools for interacting with the OpenAI API. To authorize access set the OPENAI_API_KEY environment variable to your OpenAI API key.`,
	}

	var numImages int
	var size string
	var responseFormat string
	var user string

	addImageFlags := func(cmd *cobra.Command) {
		cmd.Flags().IntVarP(&numImages, "num-images", "n", 1, "number of images (optional, default: 1)")
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
			createImage(prompt, numImages, getSize(size), getResponseFormat(responseFormat), user)
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
			imageVariations(imageFile, numImages, getSize(size), getResponseFormat(responseFormat), user)
		},
	}
	addImageFlags(createImageVariationsCmd)
	imageCmd.AddCommand(createImageVariationsCmd)

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
	res, err := client.CreateImage(openai.CreateImageReq{
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
	res, err := client.CreateImageVariations(openai.CreateImageVariationsReq{
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
