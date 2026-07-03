package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// ParseFile parses a single Go file and prints identified structs and functions
func ParseFile(filePath string) error {
	// Create the file set (holds line/column information for our code)
	fset := token.NewFileSet()

	// Parse the source file into an Abstract Syntax Tree (AST)
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse file: %w", err)
	}

	fmt.Printf("--- Parsing File: %s ---\n", filePath)
	fmt.Printf("Package: %s\n\n", node.Name.Name)

	// Inspect the AST to find Structs and Functions
	ast.Inspect(node, func(n ast.Node) bool {
		// Stop if node is nil
		if n == nil {
			return true
		}

		// Look for Type Declarations (this is where Structs are defined)
		if typeDecl, ok := n.(*ast.TypeSpec); ok {
			// Check if the type is a StructType
			if _, ok := typeDecl.Type.(*ast.StructType); ok {
				pos := fset.Position(typeDecl.Pos())
				fmt.Printf("[Struct] Name: %s (Line: %d)\n", typeDecl.Name.Name, pos.Line)
			}
		}

		// Look for Function Declarations
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			pos := fset.Position(funcDecl.Pos())

			// Check if it's a method on a struct
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				// Methods look like: func (u *User) Greet()
				recvType := funcDecl.Recv.List[0].Type
				fmt.Printf("[Method] Name: %s on Receiver Type: %v (Line: %d)\n", funcDecl.Name.Name, recvType, pos.Line)
			} else {
				// Regular functions
				fmt.Printf("[Function] Name: %s (Line: %d)\n", funcDecl.Name.Name, pos.Line)
			}
		}

		return true
	})

	return nil
}
