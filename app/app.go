package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lmninja/internal/llm"
	"lmninja/internal/security"
	"lmninja/internal/sidecar"
	"lmninja/internal/storage"
	"net/http"
	goruntime "runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)



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

	go func() {
		runtime.LogInfo(a.ctx, "[Go] Starting Python sidecar in background...")
		if err := sidecar.Start(a.ctx); err != nil {
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
	connections, err := a.db.GetConnections(a.ctx)
	if err != nil {
		return nil, err
	}
	return connections, nil
}

func (a *App) SaveConnection(meta storage.ConnectionMetadata, apiKey string) ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
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
	return a.GetConnections()
}

func (a *App) DeleteConnection(id string) ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
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

func (a *App) SelectGGUFFile() (string, error) {
	if err := a.checkReady(); err != nil {
		return "", err
	}
	runtime.LogDebug(a.ctx, "[Go] SelectGGUFFile called.")
	
	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select GGUF Model File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "GGUF Models (*.gguf)",
				Pattern:     "*.gguf",
			},
		},
	})
	if err != nil {
		return "", err
	}
	
	return filePath, nil
}


func (a *App) LoadModel(connectionID string) (storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return storage.ConnectionMetadata{}, err
	}
	runtime.LogDebug(a.ctx, fmt.Sprintf("[Go] LoadModel called for ID: %s", connectionID))

	a.mu.Lock()
	defer a.mu.Unlock()

	connections, err := a.db.GetConnections(a.ctx)
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

	if targetMeta.Provider == "gguf" {
		runtime.LogInfof(a.ctx, "[Go] Requesting Python sidecar to load GGUF model: %s", targetMeta.Model)
		err := a.requestGGUFLoad(targetMeta.Model)
		if err != nil {
			a.activeLLM = nil
			runtime.LogError(a.ctx, fmt.Sprintf("[Go] Failed to load GGUF model in sidecar: %v", err))
			return storage.ConnectionMetadata{}, fmt.Errorf("sidecar failed to load model: %w", err)
		}
		runtime.LogInfo(a.ctx, "[Go] Sidecar confirmed GGUF model is loaded.")
	}

	client, err := llm.NewLLMClient(targetMeta, a.credManager)
	if err != nil {
		a.activeLLM = nil
		return storage.ConnectionMetadata{}, fmt.Errorf("failed to create LLM client for %s: %w", targetMeta.Name, err)
	}

	a.activeLLM = client
	runtime.LogInfof(a.ctx, "[Go] Successfully activated model for chat: %s", targetMeta.Name)
	return targetMeta, nil
}

func (a *App) requestGGUFLoad(modelPath string) error {
	apiURL := "http://127.0.0.1:1337/load-gguf"
	requestBody, err := json.Marshal(map[string]string{"model_path": modelPath})
	if err != nil {
		return fmt.Errorf("failed to marshal load request: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("failed to create load request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to sidecar failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sidecar returned non-200 status (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
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
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	response, err := activeClient.Query(ctx, prompt)
	if err != nil {
		return "", err
	}
	return response, nil
}