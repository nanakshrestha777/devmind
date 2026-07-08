package parser

import (
	"devmind/db"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

func ParseFile(filePath string, database *db.DB) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var currentFuncNodeID string // Track the full ID (name:line)

	ast.Inspect(node, func(n ast.Node) bool {
		// 1. Handle Function Declarations
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			pos := fset.Position(funcDecl.Pos())
			end := fset.Position(funcDecl.End())

			// Set the current ID
			currentFuncNodeID = fmt.Sprintf("%s:%d", funcDecl.Name.Name, pos.Line)

			query := `INSERT OR REPLACE INTO nodes (id, name, type, file_path, start_line, end_line) VALUES (?, ?, ?, ?, ?, ?)`
			_, err := database.Conn.Exec(query, currentFuncNodeID, funcDecl.Name.Name, "function", filePath, pos.Line, end.Line)
			if err != nil {
				fmt.Printf("Error inserting %s: %v\n", funcDecl.Name.Name, err)
			} else {
				fmt.Printf("Indexed function: %s\n", funcDecl.Name.Name)
			}
		}

		if call, ok := n.(*ast.CallExpr); ok {
			if ident, ok := call.Fun.(*ast.Ident); ok {
				if currentFuncNodeID != "" {

					fmt.Printf("Linking: %s -> calls -> %s\n", currentFuncNodeID, ident.Name)

					query := `INSERT OR IGNORE INTO edges (from_node_id, to_node_id, type) VALUES (?, ?, ?)`
					_, err := database.Conn.Exec(query, currentFuncNodeID, ident.Name, "calls")
					if err != nil {
						fmt.Printf("Error saving edge: %v\n", err)
					}
				}
			}
		}

		return true
	})
	return nil
}

func ScanRepository(rootPath string, database *db.DB) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == "data" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		if ext == ".go" {
			fmt.Printf(">>> Scanning: %s\n", path)
			return ParseFile(path, database)
		}

		return nil
	})
}
