package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (o *openAI) makeJSONRequest(ctx context.Context, uri string, req any, resp any) error {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("openai: JSON encoding error: %w", err)
	}
	body := bytes.NewBuffer(bodyBytes)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", uri, body)
	if err != nil {
		return fmt.Errorf("openai: HTTP request creation error: %w", err)
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
	httpReq, err := http.NewRequestWithContext(ctx, "POST", uri, body)
	if err != nil {
		return fmt.Errorf("openai: HTTP request creation error: %w", err)
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	return o.makeHttpRequest(httpReq, resp)
}

func (o *openAI) makeHttpRequest(httpReq *http.Request, resp any) error {
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)
	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("X-Request-ID", "openai-go-"+strconv.FormatUint(rand.Uint64(), 16))
	if o.organization != "" {
		httpReq.Header.Set("OpenAI-Organization", o.organization)
	}
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
