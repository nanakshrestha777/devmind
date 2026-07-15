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

	// MOVE THIS OUTSIDE: It must persist across the whole file inspection
	var currentFuncNodeID string

	ast.Inspect(node, func(n ast.Node) bool {
		// 1. Handle Function Declarations
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			startLine := fset.Position(funcDecl.Pos()).Line
			cleanPath := filepath.ToSlash(filePath)

			// Update the outer variable
			currentFuncNodeID = fmt.Sprintf("%s:%s:%d", cleanPath, funcDecl.Name.Name, startLine)

			query := `INSERT OR REPLACE INTO nodes (id, name, type, file_path, start_line, end_line) VALUES (?, ?, ?, ?, ?, ?)`
			database.Conn.Exec(query, currentFuncNodeID, funcDecl.Name.Name, "function", cleanPath, startLine, fset.Position(funcDecl.End()).Line)
		}

		// 2. Handle Function Calls
		if call, ok := n.(*ast.CallExpr); ok {
			var callName string
			switch fun := call.Fun.(type) {
			case *ast.Ident:
				callName = fun.Name
			case *ast.SelectorExpr:
				if x, ok := fun.X.(*ast.Ident); ok {
					callName = fmt.Sprintf("%s.%s", x.Name, fun.Sel.Name)
				} else {
					callName = fun.Sel.Name
				}
			}

			// USE THE OUTER VARIABLE: Do not re-declare it!
			if callName != "" && currentFuncNodeID != "" {
				query := `INSERT OR IGNORE INTO edges (from_node_id, to_node_id, type) VALUES (?, ?, ?)`
				_, err := database.Conn.Exec(query, currentFuncNodeID, callName, "calls")
				if err != nil {
					fmt.Printf("Error saving edge: %v\n", err)
				}
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
