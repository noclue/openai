package openai

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCheckErrResponse(t *testing.T) {
	t.Parallel()
	t.Run("status code is 200", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			StatusCode: 200,
		}
		err := checkErrResponse(resp)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	})
	t.Run("status code is 400", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			StatusCode: 400,
		}
		err := checkErrResponse(resp)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
	t.Run("text/plain content-type returns simple error", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			StatusCode: 400,
			Header: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		}
		err := checkErrResponse(resp)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
	t.Run("application/json content-type and openAI error json return OpenAIAPIErrorDetails error", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			StatusCode: 400,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(`{
				"error": {
					"code": "invalid_api_key",
					"message": "Incorrect API key provided...",
					"param": null,
					"type": "invalid_request_error"
				}
			}`)),
		}
		err := checkErrResponse(resp)
		var openAPIErr *APIError
		if ok := errors.As(err, &openAPIErr); !ok {
			t.Errorf("expected OpenAI error error but got: %#v", err)
		}
	})
	t.Run(`application/json content-type and openAI error json missing "message" returns ErrDecodingResponse error`, func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			StatusCode: 400,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(`{
				"error": {
					"code": "invalid_api_key",
					"param": null,
					"type": "invalid_request_error"
				}
			}`)),
		}
		err := checkErrResponse(resp)
		if !errors.Is(err, ErrDecodingResponse) {
			t.Errorf("expected ErrDecodingResponse error but got: %#v", err)
		}
	})
}
