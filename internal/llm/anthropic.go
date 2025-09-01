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

func (c *AnthropicClient) Query(ctx context.Context, prompt string) (string, error) {
	apiURL := "https://api.anthropic.com/v1/messages"
	requestBody, err := json.Marshal(anthropicRequest{
		Model:     c.meta.Model,
		MaxTokens: 1024,
		Messages:  []anthropicMessage{{Role: "user", Content: prompt}},
	})
	if err != nil { return "", fmt.Errorf("failed to marshal anthropic request: %w", err) }

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil { return "", fmt.Errorf("failed to create anthropic request: %w", err) }

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.client.Do(req)
	if err != nil { return "", fmt.Errorf("failed to send request to anthropic: %w", err) }
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil { return "", fmt.Errorf("failed to read anthropic response body: %w", err) }

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Type  string `json:"type"`
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			}
		}
		if err := json.Unmarshal(bodyBytes, &errorResponse); err == nil && errorResponse.Error.Message != "" {
			return "", fmt.Errorf("Anthropic API Error (%s): %s", errorResponse.Error.Type, errorResponse.Error.Message)
		}
		return "", fmt.Errorf("Anthropic API returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse anthropicResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal successful anthropic response: %w", err)
	}

	if len(apiResponse.Content) == 0 || apiResponse.Content[0].Type != "text" {
		return "", fmt.Errorf("anthropic returned no text content")
	}

	return apiResponse.Content[0].Text, nil
}

// ListModels provides a pre-defined list for Anthropic.
func (c *AnthropicClient) ListModels(ctx context.Context) ([]string, error) {
	return []string{
		"claude-3-5-sonnet-20240620",
		"claude-3-opus-20240229",
		"claude-3-sonnet-20240229",
		"claude-3-haiku-20240307",
	},
	nil
}
