package llm

import (
	"context"
	"lmninja/internal/storage"
)

// LLM is the universal interface that all language model clients must implement.
type LLM interface {
	Query(ctx context.Context, prompt string) (string, error)
	Metadata() storage.ConnectionMetadata
	ListModels(ctx context.Context) ([]string, error)
}