package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lmninja/internal/storage"
	"net/http"
)

// AnthropicClient makes requests to the Anthropic Claude API.
type AnthropicClient struct {
	apiKey string
	meta   storage.ConnectionMetadata
	client *http.Client
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []anthropicMessage `json:"messages"`
}
type anthropicContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
type anthropicResponse struct {
	Content []anthropicContentBlock `json:"content"`
	Error   struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

func NewAnthropicClient(meta storage.ConnectionMetadata, apiKey string) *AnthropicClient {
	return &AnthropicClient{apiKey: apiKey, meta: meta, client: &http.Client{}}
}

func (c *AnthropicClient) Metadata() storage.ConnectionMetadata {
	return c.meta
}

// ListModels provides a pre-defined list for Anthropic.
func (c *AnthropicClient) ListModels(ctx context.Context) ([]string, error) {
	return []string{
		"claude-3-5-sonnet-20240620",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	}, nil
}

func (c *AnthropicClient) Query(ctx context.Context, prompt string) (string, error) {
	apiURL := "https://api.anthropic.com/v1/messages"
	requestBody, err := json.Marshal(anthropicRequest{
		Model: c.meta.Model, MaxTokens: 2048, Messages: []anthropicMessage{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal anthropic request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create anthropic request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to anthropic: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read anthropic response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse anthropicResponse
		if err := json.Unmarshal(bodyBytes, &errorResponse); err == nil && errorResponse.Error.Message != "" {
			return "", fmt.Errorf("anthropic API error (%s): %s", errorResponse.Error.Type, errorResponse.Error.Message)
		}
		return "", fmt.Errorf("anthropic API returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse anthropicResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal successful anthropic response: %w", err)
	}

	for _, block := range apiResponse.Content {
		if block.Type == "text" {
			return block.Text, nil
		}
	}

	return "", fmt.Errorf("anthropic returned no text content in response")
}
