package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sql.DB
}

func InitDB(filepath string) (*DB, error) {
	conn, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := conn.Exec("PRAGMA foreign_keys=ON;"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	db := &DB{Conn: conn}
	if err := db.createSchema(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) createSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS nodes (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		file_path TEXT NOT NULL,
		start_line INTEGER,
		end_line INTEGER
	);

	CREATE TABLE IF NOT EXISTS edges (
		from_node_id TEXT NOT NULL,
		to_node_id TEXT NOT NULL,
		type TEXT NOT NULL,
		PRIMARY KEY (from_node_id, to_node_id, type),
		FOREIGN KEY (from_node_id) REFERENCES nodes(id) ON DELETE CASCADE,
		FOREIGN KEY (to_node_id) REFERENCES nodes(id) ON DELETE CASCADE
	);`

	_, err := db.Conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
