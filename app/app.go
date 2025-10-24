package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lmninja/internal/attacks"
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
	ctx           context.Context
	db            *storage.DB
	credManager   *security.CredentialManager
	activeLLM     llm.LLM
	mu            sync.Mutex
	isReady       bool
	sidecarStatus string
}

func NewApp() *App {
	return &App{
		sidecarStatus: "starting",
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	db, err := storage.NewDB(a.ctx)
	if err != nil {
		runtime.LogFatal(a.ctx, fmt.Sprintf("FATAL: Failed to init database: %v", err))
		return
	}
	a.db = db
	a.credManager = security.NewCredentialManager()
	go func() {
		if err := sidecar.Start(a.ctx); err != nil {
			a.sidecarStatus = "error"
			runtime.EventsEmit(a.ctx, "sidecar:error", fmt.Sprintf("Failed to start Python sidecar: %v", err))
		} else {
			a.sidecarStatus = "ready"
			runtime.EventsEmit(a.ctx, "sidecar:ready")
		}
	}()
	a.isReady = true
}

func (a *App) GetSidecarStatus() string { return a.sidecarStatus }

func (a *App) Shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
	if err := sidecar.Stop(ctx); err != nil {
		runtime.LogError(ctx, fmt.Sprintf("Error stopping Python sidecar: %v", err))
	}
}

func (a *App) checkReady() error {
	if !a.isReady {
		return fmt.Errorf("backend is not ready yet")
	}
	return nil
}

func (a *App) GetAppInfo() AppInfo {
	return AppInfo{Version: "0.1.0", OS: goruntime.GOOS, Arch: goruntime.GOARCH}
}

func (a *App) GetConnections() ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil { return nil, err }
	return a.db.GetConnections(a.ctx)
}

func (a *App) SaveConnection(meta storage.ConnectionMetadata, apiKey string) ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil { return nil, err }
	isCloud := meta.Provider == "openai" || meta.Provider == "gemini" || meta.Provider == "anthropic"
	if meta.ID == "" {
		meta.ID = uuid.NewString()
		meta.CreatedAt = time.Now().UTC().Format(time.RFC3339)
		if isCloud && apiKey == "" {
			return nil, fmt.Errorf("API key cannot be empty")
		}
	}
	if isCloud && apiKey != "" {
		if err := a.credManager.StoreAPIKey(meta.ID, apiKey); err != nil {
			return nil, fmt.Errorf("failed to save API key: %w", err)
		}
	}
	if err := a.db.SaveConnection(a.ctx, meta); err != nil {
		return nil, err
	}
	return a.GetConnections()
}

func (a *App) DeleteConnection(id string) ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil { return nil, err }
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

func (a *App) RenameConnection(id string, newName string) ([]storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	if newName == "" {
		return nil, fmt.Errorf("new name cannot be empty")
	}
	if err := a.db.RenameConnection(a.ctx, id, newName); err != nil {
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
	if err := a.checkReady(); err != nil { return storage.ConnectionMetadata{}, err }
	a.mu.Lock()
	defer a.mu.Unlock()
	connections, err := a.db.GetConnections(a.ctx)
	if err != nil { return storage.ConnectionMetadata{}, err }
	var targetMeta storage.ConnectionMetadata
	found := false
	for _, c := range connections {
		if c.ID == connectionID {
			targetMeta = c
			found = true
			break
		}
	}
	if !found { return storage.ConnectionMetadata{}, fmt.Errorf("connection not found") }
	if targetMeta.Provider == "gguf" {
		if err := a.requestGGUFLoad(targetMeta); err != nil {
			a.activeLLM = nil
			return storage.ConnectionMetadata{}, fmt.Errorf("sidecar failed to load model: %w", err)
		}
	}
	client, err := llm.NewLLMClient(targetMeta, a.credManager)
	if err != nil {
		a.activeLLM = nil
		return storage.ConnectionMetadata{}, fmt.Errorf("failed to create client: %w", err)
	}
	a.activeLLM = client
	return targetMeta, nil
}

func (a *App) UnloadModel(connectionID string) error {
	if err := a.checkReady(); err != nil { return err }
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.activeLLM == nil || a.activeLLM.Metadata().ID != connectionID {
		return nil
	}
	if a.activeLLM.Metadata().Provider == "gguf" {
		if err := a.requestGGUFUnload(); err != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("Failed to unload GGUF: %v", err))
		}
	}
	a.activeLLM = nil
	return nil
}

func (a *App) GetAttackCategories() ([]attacks.AttackCategory, error) {
	if err := a.checkReady(); err != nil { return nil, err }
	return attacks.GetCategories()
}

func (a *App) GetPromptsForScan(categoryIDs []string, limit int) ([]string, error) {
	if err := a.checkReady(); err != nil { return nil, err }
	return attacks.GetPrompts(categoryIDs, limit)
}

func (a *App) EvaluatePrompt(prompt string) (string, error) {
	if err := a.checkReady(); err != nil { return "", err }
	mockResult := map[string]interface{}{"verdict": "SAFE", "details": map[string]interface{}{"prompt_injection": map[string]interface{}{"label": "benign", "score": 0.01}}}
	time.Sleep(100 * time.Millisecond)
	resultBytes, _ := json.Marshal(mockResult)
	return string(resultBytes), nil
}

func (a *App) SendMessage(prompt string) (string, error) {
	if err := a.checkReady(); err != nil { return "", err }
	a.mu.Lock()
	activeClient := a.activeLLM
	a.mu.Unlock()
	if activeClient == nil { return "", fmt.Errorf("no model loaded") }
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	return activeClient.Query(ctx, prompt)
}

func (a *App) requestGGUFLoad(meta storage.ConnectionMetadata) error {
	apiURL := "http://127.0.0.1:1337/load-gguf"
	gpuLayers := -1
	if meta.GpuLayers != nil {
		gpuLayers = *meta.GpuLayers
	}
	requestBody, _ := json.Marshal(map[string]interface{}{"model_path": meta.Model, "n_gpu_layers": gpuLayers})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return fmt.Errorf("request to sidecar failed: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sidecar error (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}

func (a *App) requestGGUFUnload() error {
	apiURL := "http://127.0.0.1:1337/unload-gguf"
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", apiURL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return fmt.Errorf("request to sidecar failed: %w", err) }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sidecar error on unload (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}