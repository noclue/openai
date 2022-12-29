package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
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
}

type openAI struct {
	// APIKey is the OpenAI API key
	APIKey string
	// Client is the http client to use to make requests to the OpenAI API
	Client HttpClient
}

type openAIOption func(*openAI)

// WithHttpClient sets the http client to use to make requests to the OpenAI
// API. One can set request and response timeouts, for example.
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

func (o *openAI) makeJSONRequest(ctx context.Context, uri string, req any, resp any) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("openai: JSON encoding error: %w", err)
	}
	body := bytes.NewBuffer(bodyBytes)
	httpReq, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return fmt.Errorf("openai: HTTP request creation error: %w", err)
	}
	if ctx != nil {
		httpReq = httpReq.WithContext(ctx)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	return o.makeHttpRequest(httpReq, resp)
}

// makeMultiPartRequest makes a multipart request to the OpenAI API. It accepts a
// map of form fields and a list of files paths to upload.
func (o *openAI) makeMultiPartRequest(ctx context.Context, uri string, fields map[string]string, files map[string]string, resp any) error {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range fields {
		if err := writer.WriteField(key, val); err != nil {
			return fmt.Errorf("openai: multipart form field encoding error: %w", err)
		}
	}
	for key, file := range files {
		fileWriter, err := writer.CreateFormFile(key, filepath.Base(file))
		if err != nil {
			return fmt.Errorf("openai: multipart form file encoding error: %w", err)
		}
		fh, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("openai: multipart form file opening error: %w", err)
		}
		defer fh.Close()
		if _, err = io.Copy(fileWriter, fh); err != nil {
			return fmt.Errorf("openai: multipart form file copying error: %w", err)
		}
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("openai: multipart form closing error: %w", err)
	}
	httpReq, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return fmt.Errorf("openai: HTTP request creation error: %w", err)
	}
	if ctx != nil {
		httpReq = httpReq.WithContext(ctx)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	return o.makeHttpRequest(httpReq, resp)
}

func (o *openAI) makeHttpRequest(httpReq *http.Request, resp any) error {
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Request-ID", "openai-go-"+strconv.FormatUint(rand.Uint64(), 16))
	httpResp, err := o.Client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("openai: HTTP error: %w", err)
	}

	if err = checkErrResponse(httpResp); err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if err = checkJSONContentType(httpResp); err != nil {
		return err
	}
	responseBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return fmt.Errorf("openai: HTTP response read error: %w", err)
	}
	if err := json.Unmarshal(responseBody, resp); err != nil {
		return fmt.Errorf("openai: HTTP success response JSON decoding error: %w", err)
	}
	return nil
}
