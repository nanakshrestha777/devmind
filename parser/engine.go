package parser

import (
	"fmt"
)

// 1. The Interface
type CodeParser interface {
	Parse(filePath string) ([]Node, []Edge, error)
}

// 2. The Implementation Structs
type GoParser struct{}

func (p *GoParser) Parse(filePath string) ([]Node, []Edge, error) {
	fmt.Printf("Parsing Go file: %s\n", filePath)
	return nil, nil, nil
}

type PythonParser struct{}

func (p *PythonParser) Parse(filePath string) ([]Node, []Edge, error) {
	fmt.Printf("Parsing Python file: %s\n", filePath)
	return nil, nil, nil
}

// 3. The Factory (Only one definition!)
func GetParser(language string) CodeParser {
	switch language {
	case "go":
		return &GoParser{}
	case "python":
		return &PythonParser{}
	default:
		return nil
	}
}
