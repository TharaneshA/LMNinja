package storage

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/google/uuid"
)

type DashboardStats struct {
	ModelsScanned      int     `json:"modelsScanned"`
	TotalScans         int     `json:"totalScans"`
	VulnerabilitiesFound int     `json:"vulnerabilitiesFound"`
	OverallPassRate    float64 `json:"overallPassRate"`
}

type VulnerabilityByModel struct {
	ModelName string `json:"modelName"`
	VulnCount int    `json:"vulnCount"`
}

type ScanHistoryItem struct {
	ID                 string `json:"id"`
	TargetModelName    string `json:"targetModelName"`
	StartTime          string `json:"startTime"`
	Status             string `json:"status"`
	TotalPrompts       int    `json:"totalPrompts"`
	VulnerabilitiesFound int    `json:"vulnerabilitiesFound"`
}

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
	createConnectionsTableSQL := `
	CREATE TABLE IF NOT EXISTS connections (
		"id" TEXT NOT NULL PRIMARY KEY, "name" TEXT NOT NULL UNIQUE, "group" TEXT,
		"provider" TEXT NOT NULL, "model" TEXT NOT NULL, "created_at" TEXT NOT NULL
	);`
	if _, err := db.Exec(createConnectionsTableSQL); err != nil { return err }

	alterConnectionsTableSQL := `ALTER TABLE connections ADD COLUMN gpu_layers INTEGER;`
	_, _ = db.Exec(alterConnectionsTableSQL)

	createScansTableSQL := `
	CREATE TABLE IF NOT EXISTS scans (
		"id" TEXT NOT NULL PRIMARY KEY,
		"target_model_name" TEXT NOT NULL,
		"start_time" TEXT NOT NULL,
		"end_time" TEXT,
		"status" TEXT NOT NULL
	);`
	if _, err := db.Exec(createScansTableSQL); err != nil { return err }

	createScanResultsTableSQL := `
	CREATE TABLE IF NOT EXISTS scan_results (
		"id" INTEGER PRIMARY KEY AUTOINCREMENT,
		"scan_id" TEXT NOT NULL,
		"prompt" TEXT NOT NULL,
		"response" TEXT NOT NULL,
		"evaluation_json" TEXT NOT NULL,
		FOREIGN KEY(scan_id) REFERENCES scans(id)
	);`
	if _, err := db.Exec(createScanResultsTableSQL); err != nil { return err }

	runtime.LogInfo(ctx, "Database tables initialized/updated successfully.")
	return nil
}

func (db *DB) GetDashboardStats(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats
	// Get unique models scanned
	db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT target_model_name) FROM scans WHERE status = 'COMPLETED'").Scan(&stats.ModelsScanned)
	// Get total scans
	db.QueryRowContext(ctx, "SELECT COUNT(id) FROM scans").Scan(&stats.TotalScans)
	// Get total vulnerabilities
	db.QueryRowContext(ctx, "SELECT COUNT(id) FROM scan_results WHERE evaluation_json LIKE '%\"verdict\":\"SUCCESSFUL_ATTACK\"%'").Scan(&stats.VulnerabilitiesFound)
	
	var totalPrompts int
	db.QueryRowContext(ctx, "SELECT COUNT(id) FROM scan_results").Scan(&totalPrompts)

	if totalPrompts > 0 {
		safeCount := totalPrompts - stats.VulnerabilitiesFound
		stats.OverallPassRate = (float64(safeCount) / float64(totalPrompts)) * 100
	}
	return stats, nil
}

func (db *DB) GetVulnerabilitiesByModel(ctx context.Context) ([]VulnerabilityByModel, error) {
	query := `
		SELECT s.target_model_name, COUNT(r.id) as vul_count
		FROM scans s
		JOIN scan_results r ON s.id = r.scan_id
		WHERE r.evaluation_json LIKE '%"verdict":"SUCCESSFUL_ATTACK"%'
		GROUP BY s.target_model_name
		ORDER BY vul_count DESC;
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var results []VulnerabilityByModel
	for rows.Next() {
		var item VulnerabilityByModel
		if err := rows.Scan(&item.ModelName, &item.VulnCount); err != nil { return nil, err }
		results = append(results, item)
	}
	return results, nil
}

func (db *DB) GetScanHistory(ctx context.Context) ([]ScanHistoryItem, error) {
	query := `
		SELECT 
			s.id, s.target_model_name, s.start_time, s.status,
			(SELECT COUNT(id) FROM scan_results WHERE scan_id = s.id) as total_prompts,
			(SELECT COUNT(id) FROM scan_results WHERE scan_id = s.id AND evaluation_json LIKE '%"verdict":"SUCCESSFUL_ATTACK"%') as vulnerabilities_found
		FROM scans s
		ORDER BY s.start_time DESC;
	`
	rows, err := db.QueryContext(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()

	var history []ScanHistoryItem
	for rows.Next() {
		var item ScanHistoryItem
		if err := rows.Scan(&item.ID, &item.TargetModelName, &item.StartTime, &item.Status, &item.TotalPrompts, &item.VulnerabilitiesFound); err != nil {
			return nil, err
		}
		history = append(history, item)
	}
	return history, nil
}

func (db *DB) CreateScanRecord(ctx context.Context, targetModelName string, startTime string) (string, error) {
	scanID := uuid.NewString()
	query := `INSERT INTO scans (id, target_model_name, start_time, status) VALUES (?, ?, ?, ?)`
	_, err := db.ExecContext(ctx, query, scanID, targetModelName, startTime, "RUNNING")
	if err != nil {
		return "", fmt.Errorf("failed to create scan record: %w", err)
	}
	return scanID, nil
}

func (db *DB) SaveScanResult(ctx context.Context, scanID string, prompt string, response string, evalJSON string) error {
	query := `INSERT INTO scan_results (scan_id, prompt, response, evaluation_json) VALUES (?, ?, ?, ?)`
	_, err := db.ExecContext(ctx, query, scanID, prompt, response, evalJSON)
	if err != nil {
		return fmt.Errorf("failed to save scan result: %w", err)
	}
	return nil
}

func (db *DB) FinalizeScanRecord(ctx context.Context, scanID string, endTime string, status string) error {
	query := `UPDATE scans SET end_time = ?, status = ? WHERE id = ?`
	_, err := db.ExecContext(ctx, query, endTime, status, scanID)
	if err != nil {
		return fmt.Errorf("failed to finalize scan record: %w", err)
	}
	return nil
}

func (db *DB) GetConnections(ctx context.Context) ([]ConnectionMetadata, error) {
	query := "SELECT id, name, \"group\", provider, model, created_at, gpu_layers FROM connections ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil { return nil, err }
	defer rows.Close()

	var connections []ConnectionMetadata
	for rows.Next() {
		var conn ConnectionMetadata
		var group sql.NullString
		var gpuLayers sql.NullInt64
		if err := rows.Scan(&conn.ID, &conn.Name, &group, &conn.Provider, &conn.Model, &conn.CreatedAt, &gpuLayers); err != nil {
			return nil, err
		}
		if group.Valid { conn.Group = group.String }
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
	if err != nil { return fmt.Errorf("failed to save connection: %w", err) }
	return nil
}

func (db *DB) DeleteConnection(ctx context.Context, id string) error {
	query := `DELETE FROM connections WHERE id = ?`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil { return fmt.Errorf("failed to delete connection: %w", err) }
	return nil
}

func (db *DB) RenameConnection(ctx context.Context, id string, newName string) error {
	query := `UPDATE connections SET name = ? WHERE id = ?`
	result, err := db.ExecContext(ctx, query, newName, id)
	if err != nil { return fmt.Errorf("failed to rename connection: %w", err) }
	rowsAffected, err := result.RowsAffected()
	if err != nil { return fmt.Errorf("failed to get rows affected: %w", err) }
	if rowsAffected == 0 { return fmt.Errorf("no connection found with ID %s", id) }
	return nil
}