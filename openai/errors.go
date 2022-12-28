package openai

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ErrDecodingResponse is the error returned when the response cannot be decoded
var ErrDecodingResponse = errors.New("openai: cannot decode error response")

// OpenAIAPIErrorDetails is the error returned from the OpenAI API
type OpenAIAPIErrorDetails struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
	Type    string `json:"type"`
}

// Error returns the error message
func (e *OpenAIAPIErrorDetails) Error() string {
	return "openai: API error:" + e.Message
}

// openAIAPIError represents the JSON payload returned from the OpenAI API
type openAIAPIError struct {
	Error *OpenAIAPIErrorDetails `json:"error"`
}

// checkErrResponse unmarshals the http response body into an error.
func checkErrResponse(resp *http.Response) error {
	if resp.StatusCode >= http.StatusOK && resp.StatusCode < 300 {
		return nil
	}
	if err := checkContentType(resp); err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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

func checkContentType(resp *http.Response) error {
	contentType := resp.Header.Get("content-type")
	if contentType != "application/json" {
		return fmt.Errorf("openai: expected JSON payload but received %s", contentType)
	}
	return nil
}
