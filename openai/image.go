package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ResponseFormat is the format of the response. Can be url or b64_json
type ResponseFormat string

const (
	// Url is the url to the image. Only present if ResponseFormat is url.
	Url ResponseFormat = "url"
	// B64_json is the base64 encoded image data. Only present if
	// ResponseFormat is b64_json.
	B64_json ResponseFormat = "b64_json"
)

// Size is the size of the image to generate. Can be Small "256x256", Medium
// "512x512", or Large "1024x1024"
type Size string

const (
	// Small is "256x256"
	Small Size = "256x256"
	// Medium is "512x512"
	Medium Size = "512x512"
	// Large is "1024x1024"
	Large Size = "1024x1024"
)

// CreateImageReq is the request body for the OpenAI API to generate an image
// from a text prompt
type CreateImageReq struct {
	// Prompt is the text prompt to generate an image from. Must be less than
	// 1000 characters.
	Prompt string `json:"prompt"`
	// N is the number of images to generate. Default is 1. Max is 10.
	N *int `json:"n,omitempty"`
	// Size is the size of the image to generate. Can be Small "256x256",
	// Medium "512x512", or Large "1024x1024"
	Size Size `json:"size,omitempty"`
	// ResponseFormat is the format of the response. Can be url or b64_json
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
	User           string         `json:"user,omitempty"`
}

// CreateImageRespData is the image data. If ResponseFormat is url, this is the
// url to the image. If ResponseFormat is b64_json, this is the base64 encoded
// image data.
type CreateImageRespData struct {
	// URL is the url to the image. Only present if ResponseFormat is url.
	URL string `json:"url,omitempty"`
	// B64JSON is the base64 encoded image data. Only present if ResponseFormat
	// is b64_json.
	B64JSON string `json:"b64_json,omitempty"`
}

// CreateImageResp is the response body from the OpenAI API to generate an
// image from a text prompt.
type CreateImageResp struct {
	// Created is the time the image was created. Unix timestamp in seconds.
	Created int64 `json:"created"`
	// Data is the image data. If ResponseFormat is url, this is the url to the
	// image. If ResponseFormat is b64_json, this is the base64 encoded image
	// data.
	Data []CreateImageRespData `json:"data"`
}

// CreateImage makes a request to the OpenAI API to generate an image from a
// text prompt.
func (o *openAI) CreateImage(req CreateImageReq) (*CreateImageResp, error) {
	// Make http request to OpenAI API to generate image from a text prompt
	// Return response body
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("openai: JSON encoding error: %w", err)
	}
	body := bytes.NewBuffer(bodyBytes)
	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", body)
	if err != nil {
		return nil, fmt.Errorf("openai: HTTP request creation error: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+o.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpResp, err := o.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai: HTTP error: %w", err)
	}

	if err = checkErrResponse(httpResp); err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	// TODO: check 204 No Content, 202 Accepted etc.
	if err = checkJSONContentType(httpResp); err != nil {
		return nil, err
	}
	responseBody, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: HTTP response read error: %w", err)
	}
	resp := &CreateImageResp{}
	if err := json.Unmarshal(responseBody, resp); err != nil {
		return nil, fmt.Errorf("openai: HTTP success response JSON decoding error: %w", err)
	}
	return resp, nil
}
