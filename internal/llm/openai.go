package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lmninja/internal/storage"
	"net/http"
	"sort"    // <-- MISSING IMPORT ADDED
	"strings" // <-- MISSING IMPORT ADDED
)

// OpenAIClient makes requests to the OpenAI API.
type OpenAIClient struct {
	apiKey string
	meta   storage.ConnectionMetadata
	client *http.Client
}

// Structures for OpenAI API request and response
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

// New struct for the /v1/models endpoint response
type openAIModel struct {
	ID      string `json:"id"`
	OwnedBy string `json:"owned_by"`
}
type openAIModelListResponse struct {
	Data []openAIModel `json:"data"`
}

func NewOpenAIClient(meta storage.ConnectionMetadata, apiKey string) *OpenAIClient {
	return &OpenAIClient{apiKey: apiKey, meta: meta, client: &http.Client{}}
}

func (c *OpenAIClient) Metadata() storage.ConnectionMetadata {
	return c.meta
}

// ListModels implementation for OpenAI
func (c *OpenAIClient) ListModels(ctx context.Context) ([]string, error) {
	apiURL := "https://api.openai.com/v1/models"
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create openai models request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to openai models endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai models api returned non-200 status: %d", resp.StatusCode)
	}

	var listResp openAIModelListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode openai models response: %w", err)
	}

	var modelIDs []string
	for _, model := range listResp.Data {
		// Filter for common chat models to keep the list clean
		if strings.HasPrefix(model.ID, "gpt-") {
			modelIDs = append(modelIDs, model.ID)
		}
	}
	sort.Strings(modelIDs) // Sort alphabetically
	return modelIDs, nil
}

func (c *OpenAIClient) Query(ctx context.Context, prompt string) (string, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"
	requestBody, err := json.Marshal(openAIRequest{
		Model:    c.meta.Model,
		Messages: []openAIMessage{{Role: "user", Content: prompt}},
	})
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

	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
				Code    string `json:"code"`
			} `json:"error"`
		}
		if err := json.Unmarshal(bodyBytes, &errorResponse); err == nil && errorResponse.Error.Message != "" {
			return "", fmt.Errorf("OpenAI API Error (%s): %s", errorResponse.Error.Code, errorResponse.Error.Message)
		}
		return "", fmt.Errorf("OpenAI API returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse openAIResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal successful openai response: %w", err)
	}
	if len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}

	return apiResponse.Choices[0].Message.Content, nil
}