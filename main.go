package main

import (
	"devmind/db"
	"devmind/parser"
	"log"
)

func main() {
	database, err := db.InitDB("devmind.db")
	if err != nil {
		log.Fatalf("Database initializatio%v", err)
	}

	defer database.Close()
	log.Println("SQLite Database initialized successfully.")

	err = parser.ParseFile("testdata/dummy.go")
	if err != nil {
		log.Fatalf("parser failed: %v", err)
	}

}
