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
	"path/filepath"
	goruntime "runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type FolderSelection struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

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
	ctx           context.Context
	db            *storage.DB
	credManager   *security.CredentialManager
	activeLLM     llm.LLM
	mu            sync.Mutex
	isReady       bool
	sidecarStatus string
}

func NewApp() *App {
	return &App{sidecarStatus: "starting"}
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

func (a *App) GetDashboardStats() (storage.DashboardStats, error) {
	if err := a.checkReady(); err != nil {
		return storage.DashboardStats{}, err
	}
	return a.db.GetDashboardStats(a.ctx)
}

func (a *App) GetVulnerabilitiesByModel() ([]storage.VulnerabilityByModel, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	return a.db.GetVulnerabilitiesByModel(a.ctx)
}

func (a *App) GetScanHistory() ([]storage.ScanHistoryItem, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	return a.db.GetScanHistory(a.ctx)
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
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	return a.db.GetConnections(a.ctx)
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

	if meta.Provider == "gguf" || meta.Provider == "ollama" || meta.Provider == "huggingface" {
		if meta.Model == "" {
			return "", fmt.Errorf("model path/name is required")
		}
		return fmt.Sprintf("Local connection for provider '%s' seems valid. Final check occurs on model load.", meta.Provider), nil
	}

	var specificClient llm.LLM
	switch meta.Provider {
	case "openai":
		specificClient = llm.NewOpenAIClient(meta, apiKey)
	case "gemini":
		specificClient = llm.NewGeminiClient(meta, apiKey)
	case "anthropic":
		specificClient = llm.NewAnthropicClient(meta, apiKey)
	default:
		return "", fmt.Errorf("test connection not implemented for provider: %s", meta.Provider)
	}

	_, err := specificClient.ListModels(a.ctx)
	if err != nil {
		return "", fmt.Errorf("connection test failed: %w", err)
	}

	return fmt.Sprintf("Successfully connected and validated credentials for %s.", meta.Name), nil
}

func (a *App) GetProviderModels(provider string, apiKey string) ([]string, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}

	meta := storage.ConnectionMetadata{
		Provider: provider,
		Model:    "temp",
	}

	var modelLister llm.LLM
	switch provider {
	case "openai":
		modelLister = llm.NewOpenAIClient(meta, apiKey)
	case "gemini":
		modelLister = llm.NewGeminiClient(meta, apiKey)
	case "anthropic":
		modelLister = llm.NewAnthropicClient(meta, apiKey)
	default:
		return nil, fmt.Errorf("model listing is not supported for local provider '%s'. Please enter the model name manually", provider)
	}

	return modelLister.ListModels(a.ctx)
}

func (a *App) SelectHuggingFaceFolder() (FolderSelection, error) {
	if err := a.checkReady(); err != nil {
		return FolderSelection{}, err
	}
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Hugging Face Model Folder",
	})
	if err != nil {
		return FolderSelection{}, err
	}
	if dirPath == "" {
		return FolderSelection{}, nil 
	}

	folderName := filepath.Base(dirPath)
	return FolderSelection{Name: folderName, Path: dirPath}, nil
}


func (a *App) SelectGGUFFile() (GGUFFile, error) {
	if err := a.checkReady(); err != nil {
		return GGUFFile{}, err
	}

	filePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Llama GGUF Model File",
		Filters: []runtime.FileFilter{
			{DisplayName: "GGUF Model Files (*.gguf)", Pattern: "*.gguf"},
		},
	})

	if err != nil {
		return GGUFFile{}, err
	}
	if filePath == "" {
		return GGUFFile{}, nil
	}

	fileName := filepath.Base(filePath)
	return GGUFFile{Name: fileName, Path: filePath}, nil
}


func (a *App) LoadModel(connectionID string) (storage.ConnectionMetadata, error) {
	if err := a.checkReady(); err != nil {
		return storage.ConnectionMetadata{}, err
	}
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
		return storage.ConnectionMetadata{}, fmt.Errorf("connection not found")
	}

	if targetMeta.Provider == "gguf" {
		if err := a.requestGGUFLoad(targetMeta); err != nil {
			a.activeLLM = nil
			return storage.ConnectionMetadata{}, fmt.Errorf("sidecar failed to load GGUF model: %w", err)
		}
	} else if targetMeta.Provider == "huggingface" {
		if err := a.requestHFLoad(targetMeta); err != nil {
			a.activeLLM = nil
			return storage.ConnectionMetadata{}, fmt.Errorf("sidecar failed to load Hugging Face model: %w", err)
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
	if err := a.checkReady(); err != nil {
		return err
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.activeLLM == nil || a.activeLLM.Metadata().ID != connectionID {
		return nil
	}


	provider := a.activeLLM.Metadata().Provider
	if provider == "gguf" || provider == "huggingface" {
		if err := a.requestGGUFUnload(); err != nil {
			runtime.LogError(a.ctx, fmt.Sprintf("Failed to unload local model (%s): %v", provider, err))
		}
	}

	a.activeLLM = nil
	return nil
}


func (a *App) GetAttackCategories() ([]attacks.AttackCategory, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	return attacks.GetCategories()
}

func (a *App) GetPromptsForScan(categoryIDs []string, limit int) ([]string, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	return attacks.GetPrompts(categoryIDs, limit)
}

func (a *App) EvaluateCompliance(prompt string, response string) (string, error) {
	if err := a.checkReady(); err != nil {
		return "", err
	}
	apiURL := "http://127.0.0.1:1337/evaluate_compliance"
	requestBody, _ := json.Marshal(map[string]string{"prompt": prompt, "response": response})
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to compliance evaluator failed: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("compliance evaluator error (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	return string(bodyBytes), nil
}

func (a *App) CreateScanRecord(targetModelName string) (string, error) {
	if err := a.checkReady(); err != nil {
		return "", err
	}
	startTime := time.Now().UTC().Format(time.RFC3339)
	return a.db.CreateScanRecord(a.ctx, targetModelName, startTime)
}

func (a *App) SaveScanResult(scanID, prompt, response, evalJSON string) error {
	if err := a.checkReady(); err != nil {
		return err
	}
	return a.db.SaveScanResult(a.ctx, scanID, prompt, response, evalJSON)
}

func (a *App) FinalizeScan(scanID, status string) error {
	if err := a.checkReady(); err != nil {
		return err
	}
	endTime := time.Now().UTC().Format(time.RFC3339)
	return a.db.FinalizeScanRecord(a.ctx, scanID, endTime, status)
}

func (a *App) GetScanResults(scanID string) ([]storage.ScanResultItem, error) {
	if err := a.checkReady(); err != nil {
		return nil, err
	}
	return a.db.GetScanResultsForScan(a.ctx, scanID)
}

func (a *App) SendMessage(prompt string) (string, error) {
	if err := a.checkReady(); err != nil {
		return "", err
	}
	a.mu.Lock()
	activeClient := a.activeLLM
	a.mu.Unlock()
	if activeClient == nil {
		return "", fmt.Errorf("no model loaded")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	return activeClient.Query(ctx, prompt)
}

func (a *App) requestHFLoad(meta storage.ConnectionMetadata) error {
	apiURL := "http://127.0.0.1:1337/load-hf-model"
	requestBody, _ := json.Marshal(map[string]interface{}{"model_path": meta.Model})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute) 
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to sidecar for HF load failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sidecar error on HF load (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
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
	if err != nil {
		return fmt.Errorf("request to sidecar failed: %w", err)
	}
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
	if err != nil {
		return fmt.Errorf("request to sidecar failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("sidecar error on unload (%d): %s", resp.StatusCode, string(bodyBytes))
	}
	return nil
}