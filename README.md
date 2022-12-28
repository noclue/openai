# openai

Welcome to the openai project! This project uses the OpenAI API to generate images based on creative prompts. The code for this project was largely written by AI using GitHub Copilot and ChatGPT.

## Getting Started

To get started, follow the instructions in the [OpenAI API guide](https://beta.openai.com/docs/quickstart) to obtain an API key and set it in the `OPENAI_API_KEY` environment variable. Then, clone this repository.

Run the following command to generate an image based on a creative prompt:

```
go run cmd/openai/openai.go help image image create "Pretty woman walking down the street."
```

This will generate an image based on the provided prompt.

To generate image variations, run the following command:

```bash
$ go run cmd/openai/openai.go image variations openai/testdata/image.png -n 2
```

To make image edits, run the following command:

```bash
$ go run cmd/openai/openai.go image edits openai/testdata/image.png "A winter forest with a winding path." -m openai/testdata/mask.png -n 2
```

## Dependencies

This project depends on the OpenAI API.

## Contributing

We welcome contributions to this project. If you have an idea for a new feature or improvement, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.