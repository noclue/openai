package openai

import (
	"context"
	"fmt"
	"net/http"
)

var modelsPath = fmt.Sprintf("%v/%v/models", basePath, apiVersion)

// ModelsResponse is response of the OpenAI API for the models endpoint
type ModelsResponse struct {
	// Data is the list of models
	Data []Model `json:"data"`
	// Object is the response object type. Should be set to "list"
	Object string `json:"object"`
}

// Model is a model
type Model struct {
	// ID is the model ID
	ID string `json:"id"`
	// Object is the model object type. Should be set to "model"
	Object string `json:"object"`
	// Created is the model creation date
	Created int64 `json:"created"`
	// OwnedBy is the organization that owns the model
	OwnedBy string `json:"owned_by"`
	// Permission is the model permissions
	Permission []any `json:"permission"`
	// Root is the model root
	Root string `json:"root"`
	// Parent is the model parent
	Parent string `json:"parent"`
}

// Models returns the list of models available to the user from the OpenAI API
func (c *openAI) Models(ctx context.Context) (*ModelsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, modelsPath, nil)
	if err != nil {
		return nil, err
	}

	var resp ModelsResponse
	if err := c.makeHttpRequest(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
