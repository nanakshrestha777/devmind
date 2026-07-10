package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type DB struct {
	Conn *sql.DB
}

func InitDB(dbpath string) (*DB, error) {
	// 1. Create the directory if it doesn't exist (e.g., 'data/')
	dir := filepath.Dir(dbpath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// 2. Open the database
	conn, err := sql.Open("sqlite", dbpath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// 3. Enable Foreign Keys
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
		PRIMARY KEY (from_node_id, to_node_id, type)
	);`

	_, err := db.Conn.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func (db *DB) GetFunctionReport() {
	query := `
	SELECT n.name, e.to_node_id 
	FROM nodes n
	LEFT JOIN edges e ON n.id = e.from_node_id
	WHERE n.type = 'function'`

	rows, err := db.Conn.Query(query)
	if err != nil {
		fmt.Printf("Error querying report: %v\n", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n--- Project Analysis Report ---")
	for rows.Next() {
		var name string
		var calls sql.NullString

		err := rows.Scan(&name, &calls)
		if err != nil {
			fmt.Printf("Error scanning row: %v\n", err)
			continue
		}

		if !calls.Valid {
			fmt.Printf("Function: %s | Calls: None\n", name)
		} else {
			fmt.Printf("Function: %s | Calls: %s\n", name, calls.String)
		}
	}
}

func (db *DB) GetImpact(functionName string) {

	query := `
	SELECT e.from_node_id 
	FROM edges e
	\WHERE to_node_id = ? AND type = 'calls'`

	rows, err := db.Conn.Query(query, functionName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return

	}

	defer rows.Close()

	fmt.Println("Impact Report for function:", functionName)

	found := false
	for rows.Next() {
		var caller string
		rows.Scan(&caller)
		fmt.Printf("Warning: Changing '%s' wll affect: %s\n", functionName, caller)
		found = true
	}
	if !found {
		fmt.Println("No dependencies found. Safe to modify")

	}
}

func (db *DB) Close() error {
	return db.Conn.Close()
}
