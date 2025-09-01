package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	// The blank import is necessary for the driver to register itself.
	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ConnectionMetadata holds information about a saved LLM connection.
// This struct will be used across the application.
type ConnectionMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Provider  string `json:"provider"` // e.g., "openai", "gemini", "anthropic", "ollama", "gguf"
	Model     string `json:"model"`    // e.g., "gpt-4o", "llama3", "/path/to/model.gguf"
	CreatedAt string `json:"createdAt"`
}

// DB is a wrapper around the sql.DB connection pool.
type DB struct {
	*sql.DB
}

// NewDB creates and initializes a new SQLite database connection.
func NewDB(ctx context.Context) (*DB, error) {
	appDataDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user config dir: %w", err)
	}
	dbDir := filepath.Join(appDataDir, "LMNinja")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create app config dir: %w", err)
	}
	dbPath := filepath.Join(dbDir, "lmninja.db")

	runtime.LogInfof(ctx, "Database path: %s", dbPath)

	// Open the database file, creating it if it doesn\'t exist.
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db := &DB{sqlDB}

	// Create tables if they don\'t exist.
	if err = db.createTables(ctx); err != nil {
		return nil, fmt.Errorf("failed to create database tables: %w", err)
	}

	return db, nil
}

// createTables executes the initial SQL schema setup.
func (db *DB) createTables(ctx context.Context) error {
	// We use TEXT for IDs (UUIDs), Provider, and Timestamps.
	// NOT NULL constraints ensure data integrity.
	createConnectionsTableSQL := `
	CREATE TABLE IF NOT EXISTS connections (
		"id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL,
		"provider" TEXT NOT NULL,
		"model" TEXT NOT NULL,
		"created_at" TEXT NOT NULL
	);`

	statement, err := db.Prepare(createConnectionsTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}

	runtime.LogInfo(ctx, "Database tables initialized successfully.")
	// Later, we will add tables for scan_results, schedules, etc. here.
	return nil
}

// --- CRUD Operations for Connections ---

// GetConnections retrieves all saved connections from the database.
func (db *DB) GetConnections(ctx context.Context) ([]ConnectionMetadata, error) {
	rows, err := db.Query("SELECT id, name, provider, model, created_at FROM connections ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []ConnectionMetadata
	for rows.Next() {
		var conn ConnectionMetadata
		if err := rows.Scan(&conn.ID, &conn.Name, &conn.Provider, &conn.Model, &conn.CreatedAt); err != nil {
			return nil, err
		}
		connections = append(connections, conn)
	}
	return connections, nil
}

// SaveConnection inserts or updates a connection in the database.
func (db *DB) SaveConnection(ctx context.Context, conn ConnectionMetadata) error {
	// "REPLACE INTO" is an SQLite-specific command that is very useful here.
	// It tries to INSERT, but if a row with the same PRIMARY KEY (id) already exists,
	// it performs an UPDATE instead. This handles both creating and editing in one step.
	query := `REPLACE INTO connections (id, name, provider, model, created_at) VALUES (?, ?, ?, ?, ?)`
	
	_, err := db.ExecContext(ctx, query, conn.ID, conn.Name, conn.Provider, conn.Model, conn.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save connection to database: %w", err)
	}
	return nil
}

// DeleteConnection removes a connection from the database by its ID.
func (db *DB) DeleteConnection(ctx context.Context, id string) error {
	query := `DELETE FROM connections WHERE id = ?`
	
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete connection from database: %w", err)
	}
	return nil
}