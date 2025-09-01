package llm

import (
	"context"
	"lmninja/internal/config"
	"net/http"
)

type AnthropicClient struct {
	apiKey string
	meta   config.ConnectionMetadata
	client *http.Client
}

func NewAnthropicClient(meta config.ConnectionMetadata, apiKey string) *AnthropicClient {
	return &AnthropicClient{apiKey: apiKey, meta: meta, client: &http.Client{}}
}

func (c *AnthropicClient) Query(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement the HTTP POST request to `https://api.anthropic.com/v1/messages`  
	// Use "x-api-key: <apiKey>" and "anthropic-version: 2023-06-01" headers.
	// Body should be JSON: {"model": c.meta.Model, "max_tokens": 1024, "messages": [{"role": "user", "content": prompt}]}
    // For now, we simulate a successful call for the TestConnection feature.
	if prompt == "Hello!" {
		return "Placeholder success response from Anthropic.", nil
	}
	return "Placeholder response from Anthropic for prompt: " + prompt, nil
}

func (c *AnthropicClient) Metadata() config.ConnectionMetadata {
	return c.meta
}