package llm

import (
	"context"
	"fmt"
	"lmninja/internal/config"
	"net/http"
)

type GeminiClient struct {
	apiKey string
	meta   config.ConnectionMetadata
	client *http.Client
}

func NewGeminiClient(meta config.ConnectionMetadata, apiKey string) *GeminiClient {
	return &GeminiClient{apiKey: apiKey, meta: meta, client: &http.Client{}}
}

func (c *GeminiClient) Query(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement the HTTP POST request to the Google Gemini API.
	// The URL is like: `https://generativelanguage.googleapis.com/v1beta/models/<model>:generateContent?key=<apiKey>`  
	// Body should be JSON: {"contents":[{"parts":[{"text": prompt}]}]}
	endpoint := fmt.Sprintf(" `https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s` ", c.meta.Model, c.apiKey)
	// For now, we simulate a successful call for the TestConnection feature.
	if prompt == "Hello!" {
		return "Placeholder success response from Gemini.", nil
	}
	return "Placeholder response from Gemini from endpoint: " + endpoint, nil
}

func (c *GeminiClient) Metadata() config.ConnectionMetadata {
	return c.meta
}