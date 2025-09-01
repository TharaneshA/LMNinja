package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// ProviderType defines the type of LLM connection.
type ProviderType string

const (
	ProviderOpenAI    ProviderType = "openai"
	ProviderGemini    ProviderType = "gemini"
	ProviderAnthropic ProviderType = "anthropic"
	// Note: We can add ProviderOllama here later
)

// ConnectionMetadata holds non-sensitive information about a saved LLM connection.
type ConnectionMetadata struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Provider  ProviderType `json:"provider"`
	Model     string       `json:"model"` // e.g., "gpt-4o", "gemini-1.5-pro", "claude-3-5-sonnet-20240620"
	CreatedAt string       `json:"createdAt"`
}

// ConfigManager handles reading and writing the connections.json file.
type ConfigManager struct {
	mu         sync.Mutex
	configPath string
}

// NewConfigManager creates a new manager.
func NewConfigManager() (*ConfigManager, error) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(appDataDir, "LMNinja")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}
	configPath := filepath.Join(configDir, "connections.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := os.WriteFile(configPath, []byte("[]"), 0644); err != nil {
			return nil, err
		}
	}
	return &ConfigManager{configPath: configPath}, nil
}

// LoadConnections reads all saved connections.
func (cm *ConfigManager) LoadConnections() ([]ConnectionMetadata, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, err
	}
	var connections []ConnectionMetadata
	if err := json.Unmarshal(data, &connections); err != nil {
		return nil, err
	}
	return connections, nil
}

// SaveConnections writes a list of connections.
func (cm *ConfigManager) SaveConnections(connections []ConnectionMetadata) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	data, err := json.MarshalIndent(connections, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cm.configPath, data, 0644)
}