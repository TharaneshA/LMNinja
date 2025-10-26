package sidecar

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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

	pythonExecutable := "python"

	scriptPath := "python-engine/main.py"

	runtime.LogInfof(ctx, "Attempting to start Python sidecar with command: %s %s", pythonExecutable, scriptPath)

	cmd = exec.CommandContext(ctx, pythonExecutable, scriptPath)
	cmd.Stdout = os.Stdout // Pipe stdout to the main app's console for visibility
	cmd.Stderr = os.Stderr // Pipe stderr for debugging

	if err := cmd.Start(); err != nil {
		runtime.LogErrorf(ctx, "Failed to start Python sidecar: %v", err)
		return fmt.Errorf("failed to start Python sidecar: %w", err)
	}

	runtime.LogInfof(ctx, "Python sidecar process started with PID: %d. Waiting for it to become healthy...", cmd.Process.Pid)

	return waitForHealthy(ctx)
}

func waitForHealthy(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	// Give the server a generous timeout to start up.
	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	for {
		select {
		case <-timeoutCtx.Done():
			_ = Stop(ctx) // Attempt to clean up the process if it's running
			runtime.LogError(ctx, "Python sidecar health check timed out after 120 seconds")
			return fmt.Errorf("python sidecar health check timed out")
		case <-ticker.C:
			// Use a short timeout for each individual request
			req, err := http.NewRequestWithContext(timeoutCtx, "GET", healthCheckURL, nil)
			if err != nil {
				continue // Can't even create a request, try again
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

// Stop terminates the Python sidecar service.
func Stop(ctx context.Context) error {
	if cmd == nil || cmd.Process == nil {
		runtime.LogInfo(ctx, "Python sidecar process not running, nothing to stop.")
		return nil
	}

	runtime.LogInfof(ctx, "Stopping Python sidecar process with PID: %d", cmd.Process.Pid)

	// Killing the process is the most reliable way to ensure it and any of its children terminate.
	if err := cmd.Process.Kill(); err != nil {
		// If the process has already exited, this will return an error, which is fine.
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
