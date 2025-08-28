package llm

import "context"

// LLM is the universal interface for any large language model.
type LLM interface {
	Name() string
	Query(ctx context.Context, prompt string) (string, error)
}