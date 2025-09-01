package llm

import (
	"fmt"
	"lmninja/internal/security"
	"lmninja/internal/storage"
)

// NewLLMClient is a factory function that creates the correct LLM client
// based on the provided connection metadata.
func NewLLMClient(meta storage.ConnectionMetadata, credManager *security.CredentialManager) (LLM, error) {
	apiKey, err := credManager.GetAPIKey(meta.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key for %s: %w", meta.Name, err)
	}

	switch meta.Provider {
	case "openai":
		return NewOpenAIClient(meta, apiKey), nil
	case "gemini":
		return NewGeminiClient(meta, apiKey), nil
	case "anthropic":
		return NewAnthropicClient(meta, apiKey), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", meta.Provider)
	}
}
