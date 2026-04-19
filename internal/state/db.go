package state

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sonic-family/code-vault/types"
)

type DB struct {
	db *sql.DB
}

func Open(path string) (*DB, error) {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state directory: %w", err)
	}

	// Open SQLite database
	sqlDB, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{db: sqlDB}

	// Initialize schema
	if err := db.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	log.Printf("Opened state database at %s", path)
	return db, nil
}

func (db *DB) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS installations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		version TEXT NOT NULL,
		installed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		is_running BOOLEAN DEFAULT FALSE
	);
	
	CREATE TABLE IF NOT EXISTS state_transitions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		installation_id INTEGER,
		from_state TEXT NOT NULL,
		to_state TEXT NOT NULL,
		transitioned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (installation_id) REFERENCES installations (id)
	);
	`

	_, err := db.db.Exec(schema)
	return err
}

func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

func (db *DB) GetInstallation(name string) (*types.InstallState, error) {
	var state types.InstallState
	state.Name = name

	query := `SELECT version, installed_at > 0 as installed, is_running FROM installations WHERE name = ?`
	err := db.db.QueryRow(query, name).Scan(&state.Version, &state.Installed, &state.Running)
	if err != nil {
		if err == sql.ErrNoRows {
			// Game not installed
			state.Installed = false
			state.Running = false
			return &state, nil
		}
		return nil, fmt.Errorf("failed to query installation: %w", err)
	}

	return &state, nil
}

func (db *DB) SetInstalled(name, version string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert or update installation
	query := `
	INSERT INTO installations (name, version, is_running)
	VALUES (?, ?, FALSE)
	ON CONFLICT(name) DO UPDATE SET
		version = excluded.version,
		installed_at = CURRENT_TIMESTAMP
	`

	_, err = tx.Exec(query, name, version)
	if err != nil {
		return fmt.Errorf("failed to set installed: %w", err)
	}

	// Record state transition
	_, err = tx.Exec(
		"INSERT INTO state_transitions (installation_id, from_state, to_state) VALUES ((SELECT id FROM installations WHERE name = ?), 'not_installed', 'installed')",
		name,
	)
	if err != nil {
		return fmt.Errorf("failed to record transition: %w", err)
	}

	return tx.Commit()
}

func (db *DB) SetRunning(name string, running bool) error {
	query := `UPDATE installations SET is_running = ? WHERE name = ?`
	_, err := db.db.Exec(query, running, name)
	if err != nil {
		return fmt.Errorf("failed to update running state: %w", err)
	}

	// Record state transition
	fromState := "stopped"
	toState := "running"
	if !running {
		fromState, toState = "running", "stopped"
	}

	_, err = db.db.Exec(
		"INSERT INTO state_transitions (installation_id, from_state, to_state) VALUES ((SELECT id FROM installations WHERE name = ?), ?, ?)",
		name, fromState, toState,
	)
	if err != nil {
		return fmt.Errorf("failed to record transition: %w", err)
	}

	return nil
}

func (db *DB) Remove(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Record state transition
	var fromState string
	err = tx.QueryRow("SELECT CASE WHEN is_running THEN 'running' ELSE 'installed' END FROM installations WHERE name = ?", name).Scan(&fromState)
	if err != nil {
		if err == sql.ErrNoRows {
			// Game not found, nothing to remove
			return nil
		}
		return fmt.Errorf("failed to determine current state: %w", err)
	}

	_, err = tx.Exec(
		"INSERT INTO state_transitions (installation_id, from_state, to_state) VALUES ((SELECT id FROM installations WHERE name = ?), ?, 'uninstalled')",
		name, fromState,
	)
	if err != nil {
		return fmt.Errorf("failed to record transition: %w", err)
	}

	// Remove installation
	_, err = tx.Exec("DELETE FROM installations WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("failed to remove installation: %w", err)
	}

	return tx.Commit()
}

func (db *DB) ListInstallations() ([]types.InstallState, error) {
	rows, err := db.db.Query("SELECT name, version, installed_at > 0 as installed, is_running FROM installations ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query installations: %w", err)
	}
	defer rows.Close()

	var installations []types.InstallState
	for rows.Next() {
		var state types.InstallState
		if err := rows.Scan(&state.Name, &state.Version, &state.Installed, &state.Running); err != nil {
			return nil, fmt.Errorf("failed to scan installation: %w", err)
		}
		installations = append(installations, state)
	}

	return installations, nil
}

func GetDefaultDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./sonic-state.db"
	}
	return filepath.Join(homeDir, ".sonic", "state.db")
}
