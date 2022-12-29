package openai

import (
	"context"
	"fmt"
)

var createEditPath = fmt.Sprintf("%v/%v/edits", basePath, apiVersion)

// EditRequest is the request to create an edit.
type EditRequest struct {
	// ID of the model to use.
	Model string `json:"model"`
	// The input text to use as a starting point for the edit.
	Input string `json:"input"`
	// The instruction that tells the model how to edit the input.
	Instruction string `json:"instruction"`
	// How many edits to generate. Defaults to 1.
	N *int `json:"n,omitempty"`
	// What sampling temperature to use. Higher values means the model will
	// take more risks. Try 0.9 for more creative applications, and 0 for ones
	// with a well-defined answer.
	Temperature *float64 `json:"temperature,omitempty"`
	// An alternative to sampling with temperature, called nucleus sampling,
	// where the model considers the results of the tokens with top_p probability
	// mass. So 0.1 means only the tokens comprising the top 10% probability mass
	// are considered.
	TopP *float64 `json:"top_p,omitempty"`
}

type EditChoice struct {
	Text  string `json:"text"`
	Index int    `json:"index"`
}

// EditResponse is the response from creating an edit.
type EditResponse struct {
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Choices []EditChoice `json:"choices"`
	Usage   Usage        `json:"usage"`
}

// Edit creates an edit. Given a prompt and an instruction, the model
// will return an edited version of the prompt.
func (c *openAI) Edit(ctx context.Context, req EditRequest) (*EditResponse, error) {
	var resp EditResponse
	err := c.makeJSONRequest(ctx, createEditPath, req, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
