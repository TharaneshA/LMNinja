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

const localAIServerURL = "http://127.0.0.1:1337"

// LocalClient communicates with the Python sidecar service.
type LocalClient struct {
	meta   storage.ConnectionMetadata
	client *http.Client
}

// Structures for the Python service API
type localQueryRequest struct {
	ModelPath   string `json:"model_path,omitempty"`
	OllamaModel string `json:"ollama_model,omitempty"`
	Prompt      string `json:"prompt"`
}

type localQueryResponse struct {
	Response string `json:"response"`
	Detail   string `json:"detail"` // Used for error messages from FastAPI
}

func NewLocalClient(meta storage.ConnectionMetadata) *LocalClient {
	return &LocalClient{
		meta:   meta,
		client: &http.Client{},
	}
}

func (c *LocalClient) Metadata() storage.ConnectionMetadata {
	return c.meta
}

// ListModels for local clients is not applicable as models are discovered
// by scanning a directory, which will be handled in the frontend logic.
// We return an empty list to satisfy the interface.
func (c *LocalClient) ListModels(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

// Query sends a prompt to the local Python AI engine.
func (c *LocalClient) Query(ctx context.Context, prompt string) (string, error) {
	apiURL := localAIServerURL + "/query-local"

	reqBody := localQueryRequest{Prompt: prompt}
	
	// Determine if it's a GGUF or Ollama model based on the provider string
	if c.meta.Provider == "gguf" {
		reqBody.ModelPath = c.meta.Model
	} else if c.meta.Provider == "ollama" {
		reqBody.OllamaModel = c.meta.Model
	} else {
		return "", fmt.Errorf("unsupported local provider: %s", c.meta.Provider)
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal local query request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create local query request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request to local AI engine: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response from local AI engine: %w", err)
	}

	var localResp localQueryResponse
	if err := json.Unmarshal(bodyBytes, &localResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response from local AI engine: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("local AI engine error (%d): %s", resp.StatusCode, localResp.Detail)
	}

	return localResp.Response, nil
}