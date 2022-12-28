package openai

import (
	"net/http"
)

// HttpClient is the interface for the http client to use to make requests to
// the OpenAI API
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// OpenAI is the OpenAI API client
type OpenAI interface {
	// CreateImage generates an image from a text prompt
	CreateImage(req CreateImageReq) (*CreateImageResp, error)
}

type openAI struct {
	// APIKey is the OpenAI API key
	APIKey string
	// Client is the http client to use to make requests to the OpenAI API
	Client HttpClient
}

type openAIOption func(*openAI)

// WithHttpClient sets the http client to use to make requests to the OpenAI API
func WithHttpClient(client HttpClient) openAIOption {
	return func(o *openAI) {
		o.Client = client
	}
}

// NewOpenAI creates a new OpenAI API client
func NewOpenAI(apiKey string, options ...openAIOption) OpenAI {
	res := &openAI{
		APIKey: apiKey,
		Client: http.DefaultClient,
	}
	for _, option := range options {
		option(res)
	}
	return res
}
