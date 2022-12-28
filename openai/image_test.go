package openai_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/noclue/openai-experiments/openai"
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
	apiKey = "sk-12345"
)

type mockHttpClient struct {
	response *http.Response
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.response, nil
}

func TestCreateImage(t *testing.T) {
	var httpClient = &mockHttpClient{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(successResponse)),
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
	}
	// Create a new OpenAI struct
	o := openai.NewOpenAI(apiKey, openai.WithHttpClient(httpClient))
	res, e := o.CreateImage(openai.CreateImageReq{
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
			Body:       ioutil.NopCloser(strings.NewReader(errInvalidToken)),
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		},
	}
	// Create a new OpenAI struct
	o := openai.NewOpenAI(apiKey, openai.WithHttpClient(httpClient))

	_, e := o.CreateImage(openai.CreateImageReq{
		Prompt: "This is a test",
	})
	if e == nil {
		t.Errorf("Expected error, got nil")
	}
	var openAIError *openai.OpenAIAPIErrorDetails
	if !errors.As(e, &openAIError) {
		t.Fatal("Expected OpenAI Error")
	}
	if openAIError.Code != "invalid_api_key" {
		t.Errorf("Expected error code invalid_api_key, got %#v", openAIError)
	}
}
