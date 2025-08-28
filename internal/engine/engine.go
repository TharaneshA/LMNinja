package engine

import (
	"context"
	"lmninja/internal/llm"
)

// RedTeamEngine is the core logic for running attacks.
type RedTeamEngine struct {
	targetLLM llm.LLM
}

// NewRedTeamEngine creates a new engine instance.
func NewRedTeamEngine(target llm.LLM) *RedTeamEngine {
	return &RedTeamEngine{
		targetLLM: target,
	}
}

// RunTests executes a series of attack prompts.
func (e *RedTeamEngine) RunTests(ctx context.Context, prompts []string) (map[string]string, error) {
	results := make(map[string]string)
	for _, prompt := range prompts {
		response, err := e.targetLLM.Query(ctx, prompt)
		if err != nil {
			results[prompt] = "Error: " + err.Error()
			continue
		}
		results[prompt] = response
	}
	return results, nil
}