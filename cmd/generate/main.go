package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/noclue/openai-experiments/openai"
)

func main() {
	// Define command line flags with default values
	prompt := flag.String("prompt", "", "prompt string")
	numImages := flag.Int("num-images", 1, "number of images (optional, default: 1)")
	size := flag.String("size", "medium", "size (small, medium, large) (optional, default: medium)")
	responseFormat := flag.String("response-format", "url", "response format (url, b64_json) (optional, default: url)")

	// Parse command line flags
	flag.Parse()

	if len(*prompt) < 5 {
		fmt.Println("Prompt must be at least 5 characters long")
		os.Exit(1)
	}

	// Validate size flag
	validSizes := []string{"small", "medium", "large"}
	if !contains(*size, validSizes) {
		fmt.Printf("Invalid size: %s\n", *size)
		os.Exit(1)
	}

	// Validate response format flag
	validResponseFormats := []string{"url", "b64_json"}
	if !contains(*responseFormat, validResponseFormats) {
		fmt.Printf("Invalid response format: %s\n", *responseFormat)
		os.Exit(1)
	}

	var sizeFlag openai.Size
	switch *size {
	case "small":
		sizeFlag = openai.Small
	case "medium":
		sizeFlag = openai.Medium
	case "large":
		sizeFlag = openai.Large
	default:
		panic("invalid size")
	}
	client := openai.NewOpenAI(os.Getenv("OPENAI_API_KEY"))
	res, err := client.CreateImage(openai.CreateImageReq{
		Prompt:         *prompt,
		N:              numImages,
		Size:           sizeFlag,
		ResponseFormat: openai.ResponseFormat(*responseFormat),
	})

	if err != nil {
		fmt.Printf("Encountered Error: %+v", err)
		os.Exit(1)
	}

	fmt.Printf("Response: %+v", res)
}

// contains checks if a string is in a slice of strings
func contains(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
