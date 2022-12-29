package openai_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/noclue/openai"
)

const (
	successResponse = `{
		"created": 1632632576,
		"data": [
				{
					"url": "https://oaidalleapiprodscus.blob.core.windows.net/private/blah-blah"
				}
			]
		}`
	errInvalidToken = `{
		"error": {
			"code": "invalid_api_key",
			"message": "Incorrect API key provided...",
			"param": null,
		    "type": "invalid_request_error"
		}
	}`
	// errNTooBig = `{
	// 	"error": {
	// 	  "code": null,
	// 	  "message": "15 is greater than the maximum of 10 - 'n'",
	// 	  "param": null,
	// 	  "type": "invalid_request_error"
	// 	}
	//   }`
	apiKey = "sk-12345"
)

type mockHttpClient struct {
	response         *http.Response
	requestValidator func(*http.Request)
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.requestValidator != nil {
		m.requestValidator(req)
	}
	return m.response, nil
}

func TestCreateImage(t *testing.T) {
	var httpClient = &mockHttpClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(successResponse)),
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
		requestValidator: func(req *http.Request) {
			if req.Method != http.MethodPost {
				t.Errorf("Expected POST, got %s", req.Method)
			}
			if req.URL.Path != "/v1/images/generations" {
				t.Errorf("Expected /v1/images/generations, got %s", req.URL.Path)
			}
			if req.Header.Get("Authorization") != "Bearer "+apiKey {
				t.Errorf("Expected Bearer %s, got %s", apiKey, req.Header.Get("Authorization"))
			}
			if req.Header.Get("Content-Type") != "application/json" {
				t.Errorf("Expected application/json, got %s", req.Header.Get("Content-Type"))
			}
		},
	}
	// Create a new OpenAI struct
	o := openai.NewOpenAI(apiKey, openai.WithHttpClient(httpClient))
	res, e := o.CreateImage(context.Background(), openai.CreateImageReq{
		Prompt: "This is a test",
	})
	if e != nil {
		t.Fatalf("Expected nil, got %#v", e)
	}
	if res.Created != 1632632576 {
		t.Errorf("Expected 1632632576, got %d", res.Created)
	}
}

// TestCreateImage tests the CreateImage method of the OpenAI struct.
func TestCreateImageError(t *testing.T) {
	var httpClient = &mockHttpClient{
		response: &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(errInvalidToken)),
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
	}
	// Create a new OpenAI struct
	o := openai.NewOpenAI(apiKey, openai.WithHttpClient(httpClient))

	_, e := o.CreateImage(context.Background(), openai.CreateImageReq{
		Prompt: "This is a test",
	})
	if e == nil {
		t.Errorf("Expected error, got nil")
	}
	var openAIError *openai.APIError
	if !errors.As(e, &openAIError) {
		t.Fatal("Expected OpenAI Error")
	}
	if openAIError.Code != "invalid_api_key" {
		t.Errorf("Expected error code invalid_api_key, got %#v", openAIError)
	}
}
