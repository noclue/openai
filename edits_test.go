package openai_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/noclue/openai"
)

const (
	// edits success response
	editsSuccessResponse = `{
		"created": 1632632576,
		"choices": [
				{
					"text": "blah-blah",
					"index": 0
				}
			]
		}`
)

var editsSuccessRequest = openai.EditRequest{
	Model:       "davinci",
	Input:       "blah-blah",
	Instruction: "blah-blah",
}

// TestCreateEdit tests the CreateEdit method. There is a positive test that
// validates the request is correctly serialized and validates the response is
// correctly deserialized. There is also a negative test that validates the
// error is correctly deserialized. Therre is also a negative test that
// validates exceptions are correctly handled.
func TestCreateEdit(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		// Create a new mock HTTP client
		var httpClient = &mockHttpClient{
			response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(editsSuccessResponse)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},
			requestValidator: func(req *http.Request) {
				if req.Method != http.MethodPost {
					t.Errorf("Expected POST, got %s", req.Method)
				}
				if req.URL.Path != "/v1/edits" {
					t.Errorf("Expected /v1/edits, got %s", req.URL.Path)
				}
				if req.Header.Get("Authorization") != "Bearer "+apiKey {
					t.Errorf("Expected Bearer %s, got %s", apiKey, req.Header.Get("Authorization"))
				}
				if req.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected application/json, got %s", req.Header.Get("Content-Type"))
				}
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Errorf("Expected nil, got %#v", err)
				}
				var request map[string]any
				err = json.Unmarshal(body, &request)
				if err != nil {
					t.Errorf("Expected nil, got %#v", err)
				}
				if request["model"] != editsSuccessRequest.Model {
					t.Errorf("Expected %s, got %s", editsSuccessRequest.Model, request["model"])
				}
				if request["input"] != editsSuccessRequest.Input {
					t.Errorf("Expected %s, got %s", editsSuccessRequest.Input, request["input"])
				}
				if request["instruction"] != editsSuccessRequest.Instruction {
					t.Errorf("Expected %s, got %s", editsSuccessRequest.Instruction, request["instruction"])
				}
			},
		}

		// Create a new OpenAI struct
		c := openai.NewOpenAI(apiKey, openai.WithHttpClient(httpClient))

		// Create a new edit
		edit, err := c.Edit(context.Background(), editsSuccessRequest)
		if err != nil {
			t.Errorf("Expected nil, got %#v", err)
		}

		// Validate the edit
		if edit.Created != 1632632576 {
			t.Errorf("Expected '1632632576', got %d", edit.Created)
		}
		if edit.Choices[0].Text != "blah-blah" {
			t.Errorf("Expected 'blah-blah', got %s", edit.Choices[0].Text)
		}
		if edit.Choices[0].Index != 0 {
			t.Errorf("Expected '0', got %d", edit.Choices[0].Index)
		}
	})

	t.Run("error", func(t *testing.T) {

		// Create a new mock HTTP client
		httpClient := &mockHttpClient{
			response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(errInvalidToken)),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			},

			requestValidator: func(req *http.Request) {
				if req.Method != http.MethodPost {
					t.Errorf("Expected POST, got %s", req.Method)
				}
				if req.URL.Path != "/v1/edits" {
					t.Errorf("Expected /v1/edits, got %s", req.URL.Path)
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
		c := openai.NewOpenAI(apiKey, openai.WithHttpClient(httpClient))

		// Create a new edit
		_, err := c.Edit(context.Background(), editsSuccessRequest)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "openai: API error:Incorrect API key provided..." {
			t.Errorf("Expected 'openai: API error:Incorrect API key provided...', got %s", err.Error())
		}
	})
}
