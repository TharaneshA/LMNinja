package app

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// Startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// StartRedTeamScan is a placeholder for the main engine logic
func (a *App) StartRedTeamScan(modelID string, categories []string) string {
	// In the future, this will trigger the engine and emit events.
	fmt.Printf("Starting scan on model '%s' with categories: %v\n", modelID, categories)
	return fmt.Sprintf("Scan initiated for model: %s", modelID)
}