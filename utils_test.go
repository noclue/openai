package openai

import (
	"net/http"
	"testing"
)

func TestCheckJSONContentType(t *testing.T) {
	t.Parallel()
	t.Run("content-type is json", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		}
		err := checkJSONContentType(resp)
		if err != nil {
			t.Errorf("expected no error but got %v", err)
		}
	})
	t.Run("content-type is not json", func(t *testing.T) {
		t.Parallel()
		resp := &http.Response{
			Header: http.Header{
				"Content-Type": []string{"text/plain"},
			},
		}
		err := checkJSONContentType(resp)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}
