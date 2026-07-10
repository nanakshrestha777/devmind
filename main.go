package main

import (
	"devmind/db"
	"devmind/parser"
	"fmt"
	"log"
	"os"
)

func main() {
	database, err := db.InitDB("data/devmind.db")
	if err != nil {
		log.Fatalf("DB Init failed: %v", err)
	}
	defer database.Conn.Close()

	// 1. Scan (Keep this to refresh your DB)
	_ = parser.ScanRepository("testdata/gin", database)

	// 2. Check Impact
	if len(os.Args) > 1 {
		targetFunc := os.Args[1] // Get function name from terminal
		database.GetImpact(targetFunc)
	} else {
		fmt.Println("Usage: go run main.go <functionName>")
	}
}
