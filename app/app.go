package app

import (
	"context"
	"fmt"
	"lmninja/internal/llm"
	"lmninja/internal/security"
	"lmninja/internal/sidecar"
	"lmninja/internal/storage"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type GGUFFile struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type App struct {
	ctx         context.Context
	db          *storage.DB
	credManager *security.CredentialManager
	activeLLM   llm.LLM
	mu          sync.Mutex
	isReady     bool // Safety flag to prevent calls before startup is complete
}

func NewApp() *App { return &App{} }

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogInfo(ctx, "[Go] App Startup initiated.")
	db, err := storage.NewDB(ctx)
	if err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("[Go] FATAL: Failed to init database: %v", err))
	}
	a.db = db
	runtime.LogInfo(ctx, "[Go] Database connection established.")
	a.credManager = security.NewCredentialManager()
	if err := sidecar.Start(ctx); err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("[Go] FATAL: Failed to start Python sidecar: %v", err))
	}
	a.isReady = true // Set ready flag AT THE END of startup
	runtime.LogInfo(ctx, "[Go] LMNinja application started successfully and is ready.")
}

func (a *App) Shutdown(ctx context.Context) {
	runtime.LogInfo(ctx, "Shutting down LMNinja application.")
	if a.db != nil {
		a.db.Close()
	}
	if err := sidecar.Stop(ctx); err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Error stopping Python sidecar: %v", err))
	}
}

// checkReady is a helper to prevent crashes from early frontend calls
func (a *App) checkReady() error {
	if !a.isReady {
		return fmt.Errorf("backend is not ready yet, please wait a moment")
	}
	return nil
}

func (a *App) GetConnections() ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	runtime.LogDebug(a.ctx, "[Go] GetConnections called.")
	connections, err := a.db.GetConnections(a.ctx)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("[Go] GetConnections ERROR: %v", err))
		return nil, err
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] GetConnections returning %d connections.", len(connections)))
	return connections, nil
}

func (a *App) SaveConnection(meta storage.ConnectionMetadata, apiKey string) ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] SaveConnection called for: %s", meta.Name))
	if meta.ID == "" && apiKey == "" && meta.Provider != "ollama" && meta.Provider != "gguf" {
		return nil, fmt.Errorf("API key cannot be empty for a new cloud connection")
	}
	if meta.ID == "" {
		meta.ID = uuid.NewString()
		meta.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	if apiKey != "" {
		if err := a.credManager.StoreAPIKey(meta.ID, apiKey); err != nil {
			return nil, fmt.Errorf("failed to save API key securely: %w", err)
		}
	}
	if err := a.db.SaveConnection(a.ctx, meta); err != nil {
		return nil, err
	}
	runtime.LogDebug(a.ctx, "[Go] SaveConnection successful.")
	return a.GetConnections()
}

func (a *App) DeleteConnection(id string) ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] DeleteConnection called for ID: %s", id))
	a.mu.Lock()
	if a.activeLLM != nil && a.activeLLM.Metadata().ID == id {
		a.activeLLM = nil
	}
	a.mu.Unlock()
	_ = a.credManager.DeleteAPIKey(id)
	if err := a.db.DeleteConnection(a.ctx, id); err != nil {
		return nil, err
	}
	return a.GetConnections()
}

func (a *App) TestConnection(meta storage.ConnectionMetadata, apiKey string) (string, error) {
	if err := a.checkReady(); err != nil {
		return "", err
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] TestConnection called for provider: %s", meta.Provider))
	if meta.Provider == "ollama" {
		testMeta := storage.ConnectionMetadata{Provider: "ollama", Model: meta.Model}
		testClient := llm.NewLocalClient(testMeta)
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		_, err := testClient.Query(ctx, "Hello! This is a connection test.")
		if err != nil {
			return "", fmt.Errorf("ollama connection test failed: %w", err)
		}
		return "Ollama connection successful!", nil
	}
	var testClient llm.LLM
	switch meta.Provider {
	case "openai":
		testClient = llm.NewOpenAIClient(meta, apiKey)
	case "gemini":
		testClient = llm.NewGeminiClient(meta, apiKey)
	case "anthropic":
		testClient = llm.NewAnthropicClient(meta, apiKey)
	default:
		return "", fmt.Errorf("testing is not supported for provider: %s", meta.Provider)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := testClient.Query(ctx, "Hello! This is a connection test.")
	if err != nil {
		return "", fmt.Errorf("connection test failed: %w", err)
	}
	return "Connection successful!", nil
}

func (a *App) GetProviderModels(provider string, apiKey string) ([]string, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] GetProviderModels called for: %s", provider))
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required to fetch models")
	}
	var client llm.LLM
	dummyMeta := storage.ConnectionMetadata{}
	switch provider {
	case "openai":
		client = llm.NewOpenAIClient(dummyMeta, apiKey)
	case "gemini":
		client = llm.NewGeminiClient(dummyMeta, apiKey)
	case "anthropic":
		client = llm.NewAnthropicClient(dummyMeta, apiKey)
	default:
		return nil, fmt.Errorf("model discovery is not supported for provider: %s", provider)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	return client.ListModels(ctx)
}

func (a *App) SelectGGUFFolder() ([]GGUFFile, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	runtime.LogDebug(a.ctx, "[Go] SelectGGUFFolder called.")
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{Title: "Select Folder Containing GGUF Models"})
	if err != nil {
		return nil, err
	}
	if dirPath == "" {
		return []GGUFFile{}, nil
	}
	var ggufFiles []GGUFFile
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".gguf") {
			ggufFiles = append(ggufFiles, GGUFFile{Name: file.Name(), Path: filepath.Join(dirPath, file.Name())})
		}
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] Found %d GGUF files.", len(ggufFiles)))
	return ggufFiles, nil
}

func (a *App) LoadModel(connectionID string) (storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return storage.ConnectionMetadata{}, err
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] LoadModel called for ID: %s", connectionID))
	a.mu.Lock()
	defer a.mu.Unlock()
	connections, err := a.GetConnections()
	if err != nil {
		return storage.ConnectionMetadata{}, err
	}
	var targetMeta storage.ConnectionMetadata
	found := false
	for _, c := range connections {
		if c.ID == connectionID {
			targetMeta = c
			found = true
			break
		}
	}
	if !found {
		return storage.ConnectionMetadata{}, fmt.Errorf("no connection found with ID: %s", connectionID)
	}
	client, err := llm.NewLLMClient(targetMeta, a.credManager)
	if err != nil {
		a.activeLLM = nil
		return storage.ConnectionMetadata{}, fmt.Errorf("failed to create LLM client for %s: %w", targetMeta.Name, err)
	}
	a.activeLLM = client
	runtime.LogInfof(a.ctx, "[Go] Successfully loaded model: %s", targetMeta.Name)
	return targetMeta, nil
}

func (a *App) SendMessage(prompt string) (string, error) {
	if err := a.checkReady(); err != nil {
		return "", err
	}
	a.mu.Lock()
	activeClient := a.activeLLM
	a.mu.Unlock()
	if activeClient == nil {
		return "", fmt.Errorf("no model is currently loaded")
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] SendMessage called for model '%s'", activeClient.Metadata().Name))
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	response, err := activeClient.Query(ctx, prompt)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("[Go] Error from model '%s': %v", activeClient.Metadata().Name, err))
		return "", err
	}
	return response, nil
}
