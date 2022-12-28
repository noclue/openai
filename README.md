# openai

Welcome to the Go client bindings library and CLI for the OpenAI API! With this tool, you can easily interact with the OpenAI API from your Go projects and command line.

Some of the things you can do with the CLI include:

* Generate images using the Image API
* Create variations of images
* Edit images with specified masks and prompt

## Features

* Image API support for generating images, variations, and edits
* Uses the remote OpenAI API

## Requirements

* Go 1.18 or newer
* spf13 cobra library
* `OPENAI_API_KEY`: To get an API key, follow these steps:
    * Go to the OpenAI website (https://openai.com/) and click on the "Sign Up" button in the top right corner of the page.
    * Fill out the sign up form with your name, email address, and password, and click the "Sign Up" button.
    * You will receive a confirmation email. Click on the link in the email to confirm your account.
    * Once you have confirmed your account, log in to the OpenAI website.
    * In the top right corner of the page, click on your user name, then select "API Key" from the dropdown menu.
    * Click the "Generate API Key" button.
    * Your API key will be displayed on the page. Copy the API key and use it as the value for the `OPENAI_API_KEY` environment variable.

## Examples

Create an image:

```bash
go run cmd/openai/openai.go help image image create "Pretty woman walking down the street."
```
Create image variations:

```bash
go run cmd/openai/openai.go image variations openai/testdata/image.png -n 2
```
Create image edits:

```bash
go run cmd/openai/openai.go image edits openai/testdata/image.png "A winter forest with a winding path." -m openai/testdata/mask.png -n 2
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.