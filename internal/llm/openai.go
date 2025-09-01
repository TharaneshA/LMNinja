package llm

import (
	"context"
	"lmninja/internal/config"
	"net/http"
)

type OpenAIClient struct {
	apiKey string
	meta   config.ConnectionMetadata
	client *http.Client
}

func NewOpenAIClient(meta config.ConnectionMetadata, apiKey string) *OpenAIClient {
	return &OpenAIClient{apiKey: apiKey, meta: meta, client: &http.Client{}}
}

func (c *OpenAIClient) Query(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement the HTTP POST request to `https://api.openai.com/v1/chat/completions`  
	// Use "Authorization: Bearer <apiKey>" header.
	// Body should be JSON: {"model": c.meta.Model, "messages": [{"role": "user", "content": prompt}]}
    // For now, we simulate a successful call for the TestConnection feature.
	if prompt == "Hello!" {
		return "Placeholder success response from OpenAI.", nil
	}
	return "Placeholder response from OpenAI for prompt: " + prompt, nil
}

func (c *OpenAIClient) Metadata() config.ConnectionMetadata {
	return c.meta
}