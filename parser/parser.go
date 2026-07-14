package parser

import (
	"devmind/db"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// 1. The Interface
type CodeParser interface {
	Parse(filePath string, database *db.DB) error
}

// 2. The Implementation Struct
type GoParser struct{}

// Move your original ParseFile logic into the GoParser struct
func (p *GoParser) Parse(filePath string, database *db.DB) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	var currentFuncNodeID string

	ast.Inspect(node, func(n ast.Node) bool {
		// 1. Handle Function Declarations
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			cleanPath := strings.ReplaceAll(filePath, "\\", "/")
			currentFuncNodeID = fmt.Sprintf("%s:%d", funcDecl.Name.Name, fset.Position(funcDecl.Pos()).Line)

			query := `INSERT OR REPLACE INTO nodes (id, name, type, file_path, start_line, end_line) VALUES (?, ?, ?, ?, ?, ?)`
			database.Conn.Exec(query, currentFuncNodeID, funcDecl.Name.Name, "function", cleanPath, fset.Position(funcDecl.Pos()).Line, fset.Position(funcDecl.End()).Line)
		}

		// 2. Handle Function Calls (The logic you were missing)
		if call, ok := n.(*ast.CallExpr); ok {
			var callName string

			// Detect Direct calls (myFunc()) AND Selector calls (db.Exec(), fmt.Println())
			if ident, ok := call.Fun.(*ast.Ident); ok {
				callName = ident.Name
			} else if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
				if ident, ok := selector.X.(*ast.Ident); ok {
					callName = fmt.Sprintf("%s.%s", ident.Name, selector.Sel.Name)
				}
			}

			// Save the Edge if we are inside a function
			if callName != "" && currentFuncNodeID != "" {
				query := `INSERT OR IGNORE INTO edges (from_node_id, to_node_id, type) VALUES (?, ?, ?)`
				database.Conn.Exec(query, currentFuncNodeID, callName, "calls")
			}
		}

		return true
	})
	return nil
}

// 3. The Factory
func GetParser(language string) CodeParser {
	switch language {
	case "Go":
		return &GoParser{}
	default:
		return nil
	}
}

// 4. Your Scanner remains the same
func ScanRepository(rootPath string, database *db.DB) error {
	fmt.Printf(">>> Starting scan on: %s\n", rootPath) // ADD THIS
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		// ... rest of your code
		if filepath.Ext(path) == ".go" {
			fmt.Printf(">>> Parsing: %s\n", path) // ADD THIS
			p := GetParser("Go")
			return p.Parse(path, database)
		}
		return nil
	})
}
