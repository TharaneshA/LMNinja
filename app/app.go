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
	goruntime "runtime"
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

type AppInfo struct {
	Version string `json:"version"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
}

type App struct {
	ctx         context.Context
	db          *storage.DB
	credManager *security.CredentialManager
	activeLLM   llm.LLM
	mu          sync.Mutex
	isReady     bool
}

func NewApp() *App { return &App{} }

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	runtime.LogInfo(a.ctx, "[Go] App Startup initiated.")
	
	db, err := storage.NewDB(a.ctx)
	if err != nil {
		runtime.LogFatal(a.ctx, fmt.Sprintf("[Go] FATAL: Failed to init database: %v", err))
		return
	}
	a.db = db
	runtime.LogInfo(a.ctx, "[Go] Database connection established.")
	a.credManager = security.NewCredentialManager()

	// **KEY CHANGE**: Start the sidecar in a background goroutine.
	// This makes the main Startup function return immediately.
	go func() {
		runtime.LogInfo(a.ctx, "[Go] Starting Python sidecar in background...")
		if err := sidecar.Start(a.ctx); err != nil {
			// We can't use LogFatal here as it exits the app.
			// We emit an event to the frontend to show a persistent error.
			runtime.LogError(a.ctx, fmt.Sprintf("[Go] FATAL: Failed to start Python sidecar: %v", err))
			runtime.EventsEmit(a.ctx, "sidecar:error", fmt.Sprintf("Failed to start Python sidecar: %v", err))
		} else {
			runtime.LogInfo(a.ctx, "[Go] Python sidecar started successfully.")
			runtime.EventsEmit(a.ctx, "sidecar:ready")
		}
	}()

	a.isReady = true
	runtime.LogInfo(a.ctx, "[Go] Wails startup sequence finished. UI is now active.")
}

func (a *App) Shutdown(ctx context.Context) {
	runtime.LogInfo(a.ctx, "Shutting down LMNinja application.")
	if a.db != nil {
		a.db.Close()
	}
	if err := sidecar.Stop(a.ctx); err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("Error stopping Python sidecar: %v", err))
	}
}

func (a *App) checkReady() error {
	if !a.isReady {
		return fmt.Errorf("backend is not ready yet, please wait a moment")
	}
	return nil
}

func (a *App) GetAppInfo() AppInfo {
    return AppInfo{
        Version: "0.1.0",
        OS:      goruntime.GOOS,
        Arch:    goruntime.GOARCH,
    }
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

	isCloud := meta.Provider == "openai" || meta.Provider == "gemini" || meta.Provider == "anthropic"

	if meta.ID == "" {
		meta.ID = uuid.NewString()
		meta.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		if isCloud && apiKey == "" {
			return nil, fmt.Errorf("API key cannot be empty for a new cloud connection")
		}
	}

	if isCloud && apiKey != "" {
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
	
    var testClient llm.LLM
    var err error

    if meta.Provider == "ollama" || meta.Provider == "gguf" {
        testClient = llm.NewLocalClient(meta)
    } else {
        if apiKey == "" {
            return "", fmt.Errorf("API Key is required to test a cloud connection")
        }
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
    }
    
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = testClient.Query(ctx, "Hello! This is a connection test.")
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
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	response, err := activeClient.Query(ctx, prompt)
	if err != nil {
		runtime.LogError(a.ctx, fmt.Sprintf("[Go] Error from model '%s': %v", activeClient.Metadata().Name, err))
		return "", err
	}
	return response, nil
}