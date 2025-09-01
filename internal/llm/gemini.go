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

// GeminiClient makes requests to the Google Gemini API.
type GeminiClient struct {
	apiKey string
	meta   storage.ConnectionMetadata
	client *http.Client
}

type geminiPart struct {
	Text string `json:"text"`
}
type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}
type geminiCandidate struct {
	Content geminiContent `json:"content"`
}
type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
	Error      struct {
		Message string `json:"message"`
	} `json:"error"`
}

func NewGeminiClient(meta storage.ConnectionMetadata, apiKey string) *GeminiClient {
	return &GeminiClient{apiKey: apiKey, meta: meta, client: &http.Client{}}
}

func (c *GeminiClient) Metadata() storage.ConnectionMetadata {
	return c.meta
}

// ListModels provides a pre-defined list for Gemini.
func (c *GeminiClient) ListModels(ctx context.Context) ([]string, error) {
	return []string{
		"gemini-1.5-pro-latest",
		"gemini-1.5-flash-latest",
		"gemini-1.0-pro",
	}, nil
}

func (c *GeminiClient) Query(ctx context.Context, prompt string) (string, error) {
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.meta.Model, c.apiKey)
	requestBody, err := json.Marshal(geminiRequest{
		Contents: []geminiContent{{Parts: []geminiPart{{Text: prompt}}}},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal gemini request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create gemini request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to gemini: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read gemini response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse geminiResponse
		if err := json.Unmarshal(bodyBytes, &errorResponse); err == nil && errorResponse.Error.Message != "" {
			return "", fmt.Errorf("gemini API error (%d): %s", resp.StatusCode, errorResponse.Error.Message)
		}
		return "", fmt.Errorf("gemini API returned non-200 status: %d - %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResponse geminiResponse
	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal successful gemini response: %w", err)
	}

	if len(apiResponse.Candidates) == 0 || len(apiResponse.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned no content in response")
	}

	return apiResponse.Candidates[0].Content.Parts[0].Text, nil
}
