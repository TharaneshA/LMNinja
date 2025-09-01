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

type OpenAIClient struct {
	apiKey string
	meta   storage.ConnectionMetadata
	client *http.Client
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}
type openAIChoice struct {
	Message openAIMessage `json:"message"`
}
type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
	Error   struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewOpenAIClient(meta storage.ConnectionMetadata, apiKey string) *OpenAIClient {
	return &OpenAIClient{apiKey: apiKey, meta: meta, client: &http.Client{}}
}

func (c *OpenAIClient) Metadata() storage.ConnectionMetadata {
	return c.meta
}

func (c *OpenAIClient) Query(ctx context.Context, prompt string) (string, error) {
	// ... (Query function is correct and does not need to be changed)
	apiURL := "https://api.openai.com/v1/chat/completions"
	requestBody, err := json.Marshal(openAIRequest{Model: c.meta.Model, Messages: []openAIMessage{{Role: "user", Content: prompt}}})
	if err != nil {
		return "", fmt.Errorf("failed to marshal openai request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create openai request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to openai: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read openai response body: %w", err)
	}
	var apiResponse openAIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal openai response (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	if resp.StatusCode != http.StatusOK {
		if apiResponse.Error.Message != "" {
			return "", fmt.Errorf("openai api error (%d): %s", resp.StatusCode, apiResponse.Error.Message)
		}
		return "", fmt.Errorf("openai api returned non-200 status: %d", resp.StatusCode)
	}
	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return apiResponse.Choices[0].Message.Content, nil
}
