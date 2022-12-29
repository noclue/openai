package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrDecodingResponse is the error returned when the response cannot be decoded
var ErrDecodingResponse = errors.New("openai: cannot decode error response")

// APIError is the error returned from the OpenAI API
type APIError struct {
	Code    any    `json:"code"`
	Message string `json:"message"`
	Details string `json:"param"`
	Type    string `json:"type"`
}

// Error returns the error message
func (e *APIError) Error() string {
	return "openai: API error:" + e.Message
}

// openAIAPIError represents the JSON payload returned from the OpenAI API
type openAIAPIError struct {
	Error *APIError `json:"error"`
}

// checkErrResponse unmarshals the http response body into an error.
func checkErrResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < 300 {
		return nil
	}
	if err := checkJSONContentType(resp); err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("openai: HTTP error response read error: %w", err)
	}

	openAIErr := openAIAPIError{}
	if err := json.Unmarshal(body, &openAIErr); err == nil &&
		openAIErr.Error != nil && openAIErr.Error.Message != "" {
		return openAIErr.Error
	}

	return fmt.Errorf("openai: cannot read error response with status code: %v. %w", resp.StatusCode, ErrDecodingResponse)
}
