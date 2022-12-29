package openai

import (
	"fmt"
	"net/http"
)

func checkJSONContentType(resp *http.Response) error {
	contentType := resp.Header.Get("content-type")
	if contentType != "application/json" {
		return fmt.Errorf("openai: expected JSON payload but received %s", contentType)
	}
	return nil
}
