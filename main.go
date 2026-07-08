package main

import (
	"devmind/db"
	"devmind/parser"
	"log"
)

func main() {
	database, err := db.InitDB("data/devmind.db")
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	defer database.Conn.Close()
	log.Println("Database initialized successfully.")

	err = parser.ScanRepository("testdata", database)
	if err != nil {
		log.Fatalf("parser failed: %v", err)
	}

	log.Println("Parsing completed successfully.")

	database.GetFunctionReport()

}
