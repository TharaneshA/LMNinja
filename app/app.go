package app

import (
	"context"
	"fmt"
	"lmninja/internal/config"
	"lmninja/internal/llm"
	"lmninja/internal/security"
	"lmninja/internal/sidecar" // Import the new package
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx           context.Context
	configManager *config.ConfigManager
	credManager   *security.CredentialManager
	activeLLM     llm.LLM
	mu            sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. Initialize managers here.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	cm, err := config.NewConfigManager()
	if err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("Failed to init config manager: %v", err))
		os.Exit(1)
	}
	a.configManager = cm
	a.credManager = security.NewCredentialManager()
	// Start the Python sidecar service
	if err := sidecar.Start(ctx); err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("Failed to start Python sidecar: %v", err))
		os.Exit(1)
	}
	runtime.LogInfo(ctx, "LMNinja application started successfully.")
}

// GetConnections returns the list of all saved LLM connections.
func (a *App) GetConnections() ([]config.ConnectionMetadata, error) {
	connections, err := a.configManager.LoadConnections()
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Error loading connections: %v", err))
		return nil, err
	}
	return connections, nil
}

// SaveConnection saves a new or updated connection's metadata and credentials.
func (a *App) SaveConnection(meta config.ConnectionMetadata, apiKey string) ([]config.ConnectionMetadata, error) {
	if meta.ID == "" && apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty for a new connection")
	}
	connections, err := a.configManager.LoadConnections()
	if err != nil { return nil, err }
	isUpdate := false
	if meta.ID == "" {
		meta.ID = uuid.NewString()
		meta.CreatedAt = time.Now().Format(time.RFC3339)
	}
	for i, c := range connections {
		if c.ID == meta.ID {
			connections[i] = meta
			isUpdate = true
			break
		}
	}
	if !isUpdate {
		connections = append(connections, meta)
	}
	if apiKey != "" {
		if err := a.credManager.StoreAPIKey(meta.ID, apiKey); err != nil {
			return nil, fmt.Errorf("failed to save API key securely: %w", err)
		}
	}
	if err := a.configManager.SaveConnections(connections); err != nil { return nil, err }
	return connections, nil
}

// DeleteConnection removes a connection and its associated credentials.
func (a *App) DeleteConnection(id string) ([]config.ConnectionMetadata, error) {
	a.mu.Lock()
	if a.activeLLM != nil && a.activeLLM.Metadata().ID == id {
		a.activeLLM = nil
	}
	a.mu.Unlock()
	connections, err := a.configManager.LoadConnections()
	if err != nil { return nil, err }
	var updatedConnections []config.ConnectionMetadata
	found := false
	for _, c := range connections {
		if c.ID == id {
			found = true
			continue
		}
		updatedConnections = append(updatedConnections, c)
	}
	if !found { return nil, fmt.Errorf("connection with id %s not found", id) }
	_ = a.credManager.DeleteAPIKey(id)
	if err := a.configManager.SaveConnections(updatedConnections); err != nil { return nil, err }
	return updatedConnections, nil
}

// TestConnection makes a REAL API call to verify credentials.
func (a *App) TestConnection(meta config.ConnectionMetadata, apiKey string) (string, error) {
	var testClient llm.LLM
	switch meta.Provider {
	case config.ProviderOpenAI: testClient = llm.NewOpenAIClient(meta, apiKey)
	case config.ProviderGemini: testClient = llm.NewGeminiClient(meta, apiKey)
	case config.ProviderAnthropic: testClient = llm.NewAnthropicClient(meta, apiKey)
	default: return "", fmt.Errorf("unknown provider for testing: %s", meta.Provider)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := testClient.Query(ctx, "Hello! This is a connection test.")
	if err != nil { return "", fmt.Errorf("connection test failed: %w", err) }
	return "Connection successful!", nil
}

// LoadModel finds a connection by ID, creates its client, and sets it as the active LLM.
func (a *App) LoadModel(connectionID string) (config.ConnectionMetadata, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	connections, err := a.GetConnections()
	if err != nil { return config.ConnectionMetadata{}, fmt.Errorf("could not load connections to find model: %w", err) }
	var targetMeta config.ConnectionMetadata
	found := false
	for _, c := range connections {
		if c.ID == connectionID {
			targetMeta = c
			found = true
			break
		}
	}
	if !found { return config.ConnectionMetadata{}, fmt.Errorf("no connection found with ID: %s", connectionID) }
	client, err := llm.NewLLMClient(targetMeta, a.credManager)
	if err != nil {
		a.activeLLM = nil
		return config.ConnectionMetadata{}, fmt.Errorf("failed to create LLM client for %s: %w", targetMeta.Name, err)
	}
	a.activeLLM = client
	runtime.LogInfof(a.ctx, "Successfully loaded model: %s", targetMeta.Name)
	return targetMeta, nil
}

// SendMessage sends a prompt to the currently active LLM.
func (a *App) SendMessage(prompt string) (string, error) {
	a.mu.Lock()
	activeClient := a.activeLLM
	a.mu.Unlock()
	if activeClient == nil { return "", fmt.Errorf("no model is currently loaded. Please select a model from the dropdown") }
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	runtime.LogInfof(a.ctx, "Sending prompt to active model '%s'", activeClient.Metadata().Name)
	response, err := activeClient.Query(ctx, prompt)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Error from model '%s': %v", activeClient.Metadata().Name, err))
		return "", err
	}
	return response, nil
}

// Shutdown is called when the app is closing.
func (a *App) Shutdown(ctx context.Context) {
	runtime.LogInfo(ctx, "Shutting down LMNinja application.")
	if err := sidecar.Stop(ctx); err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Error stopping Python sidecar: %v", err))
	}
}