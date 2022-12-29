package openai

import (
	"context"
	"fmt"
)

var moderationPath = fmt.Sprintf("%v/%v/moderations", basePath, apiVersion)

const (
	CategoryHate            = "hate"
	CategoryHateThreatening = "hate/threatening"
	CategorySelfHarm        = "self-harm"
	CategorySexual          = "sexual"
	CategorySexualMinors    = "sexual/minors"
	CategoryViolence        = "violence"
	CategoryViolenceGraphic = "violence/graphic"
)

// ModerationRequest is the request body for the OpenAI API moderation endpoint
type ModerationRequest struct {
	// Input is the text to be moderated
	Input []string `json:"input"`
	// Model is the model to use for moderation. Two content moderations models
	// are available: text-moderation-stable and text-moderation-latest.
	//
	// The default is text-moderation-latest which will be automatically
	// upgraded over time. This ensures you are always using our most accurate
	// model. If you use text-moderation-stable, we will provide advanced
	// notice before updating the model. Accuracy of text-moderation-stable may
	// be slightly lower than for text-moderation-latest.
	Model string `json:"model,omitempty"`
}

type ModerationResponse struct {
	// Id is the moderation ID
	ID string `json:"id"`
	// Model is the model used for moderation
	Model string `json:"model"`
	// Results is the list of moderation results
	Results []ModerationResult `json:"results"`
}

type ModerationResult struct {
	// Categories is the list of categories
	Categories map[string]bool `json:"categories"`
	// CategoryScores is the list of category scores
	CategoryScores map[string]float64 `json:"category_scores"`
	// Flagged is true if the text is flagged
	Flagged bool `json:"flagged"`
}

// Moderation returns the moderation results for the given text from the OpenAI API
func (c *openAI) Moderation(ctx context.Context, req ModerationRequest) (*ModerationResponse, error) {
	res := &ModerationResponse{}
	if err := c.makeJSONRequest(ctx, moderationPath, req, &res); err != nil {
		return nil, err
	}
	return res, nil
}
