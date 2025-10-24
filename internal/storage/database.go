package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ConnectionMetadata struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Group     string `json:"group,omitempty"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	CreatedAt string `json:"createdAt"`
	GpuLayers *int   `json:"gpuLayers,omitempty"`
}

type DB struct {
	*sql.DB
}

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
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db := &DB{sqlDB}
	if err = db.createOrUpdateTables(ctx); err != nil {
		return nil, fmt.Errorf("failed to create/update database tables: %w", err)
	}
	return db, nil
}

func (db *DB) createOrUpdateTables(ctx context.Context) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS connections (
		"id" TEXT NOT NULL PRIMARY KEY,
		"name" TEXT NOT NULL UNIQUE,
		"group" TEXT,
		"provider" TEXT NOT NULL,
		"model" TEXT NOT NULL,
		"created_at" TEXT NOT NULL
	);`
	if _, err := db.Exec(createTableSQL); err != nil {
		return err
	}

	alterTableSQL := `ALTER TABLE connections ADD COLUMN gpu_layers INTEGER;`
	_, _ = db.Exec(alterTableSQL)

	runtime.LogInfo(ctx, "Database tables initialized/updated successfully.")
	return nil
}

func (db *DB) GetConnections(ctx context.Context) ([]ConnectionMetadata, error) {
	query := "SELECT id, name, \"group\", provider, model, created_at, gpu_layers FROM connections ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []ConnectionMetadata
	for rows.Next() {
		var conn ConnectionMetadata
		var group sql.NullString
		var gpuLayers sql.NullInt64

		if err := rows.Scan(&conn.ID, &conn.Name, &group, &conn.Provider, &conn.Model, &conn.CreatedAt, &gpuLayers); err != nil {
			return nil, err
		}
		if group.Valid {
			conn.Group = group.String
		}
		if gpuLayers.Valid {
			val := int(gpuLayers.Int64)
			conn.GpuLayers = &val
		}
		connections = append(connections, conn)
	}
	return connections, nil
}

func (db *DB) SaveConnection(ctx context.Context, conn ConnectionMetadata) error {
	query := `REPLACE INTO connections (id, name, "group", provider, model, created_at, gpu_layers) VALUES (?, ?, ?, ?, ?, ?, ?)`
	
	var gpuLayers sql.NullInt64
	if conn.GpuLayers != nil {
		gpuLayers.Valid = true
		gpuLayers.Int64 = int64(*conn.GpuLayers)
	}

	_, err := db.ExecContext(ctx, query, conn.ID, conn.Name, conn.Group, conn.Provider, conn.Model, conn.CreatedAt, gpuLayers)
	if err != nil {
		return fmt.Errorf("failed to save connection to database: %w", err)
	}
	return nil
}

func (db *DB) DeleteConnection(ctx context.Context, id string) error {
	query := `DELETE FROM connections WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete connection from database: %w", err)
	}
	return nil
}

func (db *DB) RenameConnection(ctx context.Context, id string, newName string) error {
	query := `UPDATE connections SET name = ? WHERE id = ?`
	
	result, err := db.ExecContext(ctx, query, newName, id)
	if err != nil {
		return fmt.Errorf("failed to rename connection in database: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after rename: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no connection found with ID %s to rename", id)
	}
	
	return nil
}