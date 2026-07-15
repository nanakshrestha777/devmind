package main

import (
	"devmind/ai"
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

	// 2. Command routing
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go impact <functionName>")
		fmt.Println("  go run main.go explain <functionName>")
		return
	}

	command := os.Args[1]    // "impact" or "explain"
	targetFunc := os.Args[2] // The function name

	switch command {
	case "impact":
		database.GetImpact(targetFunc)

	case "explain":
		data := database.GetConnectionsForAI(targetFunc)
		fmt.Println("Asking devMind AI to analyze...")
		explanation, err := ai.Explain(data)
		if err != nil {
			log.Fatalf("AI Explanation failed: %v", err)
		}
		fmt.Println("\ndevMind AI Explanation:\n", explanation)

	default:
		fmt.Println("Unknown command. Use 'impact' or 'explain'.")
	}
}

