package openai

import (
	"strconv"
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

// ImageSize is the size of the image to generate. Can be Small "256x256", Medium
// "512x512", or Large "1024x1024"
type ImageSize string

const (
	// SmallImage is "256x256"
	SmallImage ImageSize = "256x256"
	// MediumImage is "512x512"
	MediumImage ImageSize = "512x512"
	// LargeImage is "1024x1024"
	LargeImage ImageSize = "1024x1024"
)

type CommonImageReq struct {
	// N is the number of images to generate. Default is 1. Max is 10.
	N *int `json:"n,omitempty"`
	// Size is the size of the image to generate. Can be Small "256x256",
	// Medium "512x512", or Large "1024x1024"
	Size ImageSize `json:"size,omitempty"`
	// ResponseFormat is the format of the response. Can be url or b64_json
	ResponseFormat ResponseFormat `json:"response_format,omitempty"`
	// A unique identifier representing your end-user, which can help OpenAI
	// to monitor and detect abuse. See [End User Ids] for details
	//
	// [End User Ids]: https://beta.openai.com/docs/guides/safety-best-practices/end-user-ids
	User string `json:"user,omitempty"`
}

// ImageData is the image data. If ResponseFormat is url, this is the
// url to the image. If ResponseFormat is b64_json, this is the base64 encoded
// image data.
type ImageData struct {
	// URL is the url to the image. Only present if ResponseFormat is url.
	URL string `json:"url,omitempty"`
	// B64JSON is the base64 encoded image data. Only present if ResponseFormat
	// is b64_json.
	B64JSON string `json:"b64_json,omitempty"`
}

// ImageResponse is the response body from the OpenAI API to generate an
// image from a text prompt.
type ImageResponse struct {
	// Created is the time the image was created. Unix timestamp in seconds.
	Created int64 `json:"created"`
	// Data is the image data. If ResponseFormat is url, this is the url to the
	// image. If ResponseFormat is b64_json, this is the base64 encoded image
	// data.
	Data []ImageData `json:"data"`
}

// CreateImageReq is the request body for the OpenAI API to generate an image
// from a text prompt
type CreateImageReq struct {
	CommonImageReq
	// Prompt is the text prompt to generate an image from. Must be less than
	// 1000 characters.
	Prompt string `json:"prompt"`
}

// CreateImage makes a request to the OpenAI API to generate an image from a
// text prompt.
func (o *openAI) CreateImage(req CreateImageReq) (*ImageResponse, error) {
	resp := &ImageResponse{}
	err := o.makeJSONRequest(createImagePath, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateImageVariationReq is the request body for the OpenAI API to generate
// image variations.
type CreateImageVariationsReq struct {
	CommonImageReq

	// Image is the path to an image file to generate variations from. Must be
	// a valid .png image, sqaure in shape, and less than 4MB in size.
	Image string
}

// CreateImageVariations makes a request to the OpenAI API to generate image
// variations.
func (o *openAI) CreateImageVariations(req CreateImageVariationsReq) (*ImageResponse, error) {
	resp := &ImageResponse{}
	params := map[string]string{}
	if req.N != nil {
		params["n"] = strconv.Itoa(*req.N)
	}
	if req.Size != "" {
		params["size"] = string(req.Size)
	}
	if req.ResponseFormat != "" {
		params["response_format"] = string(req.ResponseFormat)
	}
	if req.User != "" {
		params["user"] = req.User
	}
	err := o.makeMultiPartRequest(createImageVariationsPath, params, map[string]string{"image": req.Image}, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CreateImageEditsReq contains the request paramters for the OpenAI API to
// generate image edits.
type CreateImageEditsReq struct {
	CommonImageReq

	// Image is the path to an image file to generate variations from. Must be
	// a valid .png image, sqaure in shape, and less than 4MB in size.
	// (required)
	Image string
	// Mask is the path to an image file to use as a mask. Must be a valid .png
	// image, sqaure in shape, and less than 4MB in size. (optional)
	Mask string
	// Prompt is a text description of the desired image. Must be less than
	// 1000 characters. (required)
	Prompt string
}

// CreateImageEdits creates an edited or extended image given an original image
// and a prompt.
func (o *openAI) CreateImageEdits(req CreateImageEditsReq) (*ImageResponse, error) {
	resp := &ImageResponse{}
	params := map[string]string{}
	params["prompt"] = req.Prompt
	if req.N != nil {
		params["n"] = strconv.Itoa(*req.N)
	}
	if req.Size != "" {
		params["size"] = string(req.Size)
	}
	if req.ResponseFormat != "" {
		params["response_format"] = string(req.ResponseFormat)
	}
	if req.User != "" {
		params["user"] = req.User
	}
	files := map[string]string{"image": req.Image}
	if req.Mask != "" {
		files["mask"] = req.Mask
	}
	err := o.makeMultiPartRequest(createImageEditsPath, params, files, resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
