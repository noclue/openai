# openai-experiments

Welcome to the openai-experiments project! This project uses the OpenAI API to generate images based on creative prompts. The code for this project was largely written by AI using GitHub Copilot and ChatGPT.

## Getting Started

To get started, follow the instructions in the [OpenAI API guide](https://beta.openai.com/docs/quickstart) to obtain an API key and set it in the `OPENAI_API_KEY` environment variable. Then, clone this repository and run the following command:

```
go run ./cmd/generate.go --prompt "Create a Times magazine cover featuring a lunar tomato farmer with the Earth in the background. The farmer should be shown picking tomatoes on the moon, with the Earth visible in the distance. The cover should have the title "The Future of Farming: The First Tomato Farmer on the Moon" in bold, eye-catching letters. The cover should convey a sense of excitement and possibility, as the farmer represents the next step in humanity's journey to explore and thrive in the universe."
```

This will generate an image based on the provided prompt.

## Dependencies

This project depends on the OpenAI API. No external libraries are required, as the project uses the Go standard HTTP module.

## Contributing

We welcome contributions to this project. If you have an idea for a new feature or improvement, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.