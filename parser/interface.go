type CodeParser interface {
	Parse(filePath string) ([]Node, []Edge, error)
}


func GetParser (language String) codeParser {
	switch
}