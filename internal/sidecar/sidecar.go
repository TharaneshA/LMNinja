package sidecar

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	pythonServerHost = "127.0.0.1"
	pythonServerPort = "1337"
	healthCheckURL   = "http://" + pythonServerHost + ":" + pythonServerPort + "/health"
)

var cmd *exec.Cmd

func Start(ctx context.Context) error {
	ex, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exPath := filepath.Dir(ex)


	var command string
	var args []string

	sidecarExePath := filepath.Join(exPath, "lmninja-engine", "lmninja-engine.exe")

	if _, err := os.Stat(sidecarExePath); err == nil {
		command = sidecarExePath
		runtime.LogInfof(ctx, "Production mode detected. Running packaged Python sidecar.")
	} else {
		command = "python"
		args = []string{"python-engine/main.py"}
		runtime.LogInfof(ctx, "Development mode detected. Running Python script directly.")
	}
	
	runtime.LogInfof(ctx, "Attempting to start sidecar with command: %s %v", command, args)

	cmd = exec.CommandContext(ctx, command, args...)
	cmd.Stdout = os.Stdout 
	cmd.Stderr = os.Stderr 

	if err := cmd.Start(); err != nil {
		runtime.LogErrorf(ctx, "Failed to start sidecar: %v", err)
		return fmt.Errorf("failed to start sidecar: %w", err)
	}

	runtime.LogInfof(ctx, "Sidecar process started with PID: %d. Waiting for it to become healthy...", cmd.Process.Pid)

	return waitForHealthy(ctx)
}

func waitForHealthy(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			_ = Stop(ctx)
			runtime.LogError(ctx, "Python sidecar health check timed out after 120 seconds")
			return fmt.Errorf("python sidecar health check timed out")
		case <-ticker.C:
			req, err := http.NewRequestWithContext(timeoutCtx, "GET", healthCheckURL, nil)
			if err != nil {
				continue
			}

			client := &http.Client{Timeout: 200 * time.Millisecond}
			resp, err := client.Do(req)

			if err == nil && resp.StatusCode == http.StatusOK {
				runtime.LogInfo(ctx, "Python sidecar is healthy and running.")
				resp.Body.Close()
				return nil
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

func Stop(ctx context.Context) error {
	if cmd == nil || cmd.Process == nil {
		runtime.LogInfo(ctx, "Python sidecar process not running, nothing to stop.")
		return nil
	}

	runtime.LogInfof(ctx, "Stopping Python sidecar process with PID: %d", cmd.Process.Pid)

	if err := cmd.Process.Kill(); err != nil {
		if err.Error() == "os: process already finished" {
			runtime.LogWarning(ctx, "Attempted to stop Python sidecar, but process was already finished.")
			return nil
		}
		runtime.LogErrorf(ctx, "Failed to kill Python sidecar process: %v", err)
		return fmt.Errorf("failed to kill sidecar process: %w", err)
	}

	runtime.LogInfo(ctx, "Python sidecar process stopped.")
	return nil
}