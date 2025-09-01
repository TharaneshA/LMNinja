package llm

import (
	"fmt"
	"lmninja/internal/config"
	"lmninja/internal/security"
)

// NewLLMClient is a factory function that creates the correct LLM client
// based on the provided connection metadata.
func NewLLMClient(meta config.ConnectionMetadata, credManager *security.CredentialManager) (LLM, error) {
	apiKey, err := credManager.GetAPIKey(meta.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get API key for %s: %w", meta.Name, err)
	}

	switch meta.Provider {
	case config.ProviderOpenAI:
		return NewOpenAIClient(meta, apiKey), nil
	case config.ProviderGemini:
		return NewGeminiClient(meta, apiKey), nil
	case config.ProviderAnthropic:
		return NewAnthropicClient(meta, apiKey), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", meta.Provider)
	}
}