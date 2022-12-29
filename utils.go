package openai

import (
	"fmt"
	"mime"
	"net/http"
)

// checkJSONContentType checks the content-type header of the response to
// ensure it is JSON. It returns error if the content-type is not JSON or nil
// otherwise.
func checkJSONContentType(resp *http.Response) error {
	contentType := resp.Header.Get("content-type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return fmt.Errorf("openai: error parsing content-type: %v, error %w", contentType, err)
	}
	if mediaType != "application/json" {
		return fmt.Errorf("openai: content-type is not application/json: %v", contentType)
	}
	return nil
}
