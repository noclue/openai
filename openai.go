package openai

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
)

const (
	openaiGoVersion = "0.1.0"
	basePath        = "https://api.openai.com"
	apiVersion      = "v1"
)

var userAgent = fmt.Sprintf("openai-go/%v (%v; %v)", openaiGoVersion, runtime.Version(), runtime.GOOS)

// HttpClient is the interface for the http client to use to make requests to
// the OpenAI API
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// OpenAI is the OpenAI API client
type OpenAI interface {
	// CreateImage generates an image from a text prompt
	CreateImage(ctx context.Context, req CreateImageReq) (*ImageResponse, error)
	// CreateImageVariations generates image variations
	CreateImageVariations(ctx context.Context, req CreateImageVariationsReq) (*ImageResponse, error)
	// CreateImageEdits generates image edits
	CreateImageEdits(ctx context.Context, req CreateImageEditsReq) (*ImageResponse, error)
	// CreateCompletion creates a completion
	CreateCompletion(ctx context.Context, req CompletionsRequest) (*CompletionsResponse, error)
	// Edit creates an edit
	Edit(ctx context.Context, req EditRequest) (*EditResponse, error)
	// Models returns the list of models available to the user from the OpenAI API
	Models(ctx context.Context) (*ModelsResponse, error)
	// Moderation returns the moderation status of a text.
	Moderation(ctx context.Context, req ModerationRequest) (*ModerationResponse, error)
}

type openAI struct {
	// APIKey is the OpenAI API key
	APIKey string
	// Client is the http client to use to make requests to the OpenAI API
	Client HttpClient
	// organization is the organization to use for the requests to the OpenAI API.
	// See https://beta.openai.com/docs/api-reference/requesting-organization
	organization string
}

type openAIOption func(*openAI)

// WithHttpClient sets the http client to use to make requests to the OpenAI
// API. One can set request and response timeouts, for example.
func WithHttpClient(client HttpClient) openAIOption {
	return func(o *openAI) {
		o.Client = client
	}
}

// WithOrganization sets the organization to use for the requests to the OpenAI
// API. See https://beta.openai.com/docs/api-reference/requesting-organization
func WithOrganization(organization string) openAIOption {
	return func(o *openAI) {
		o.organization = organization
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
