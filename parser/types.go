package parser

type Node struct {
	Name     string
	Type     string
	Filepath string
	Line     int
}

type Edge struct {
	FromNode string
	ToNode   string
	Type     string
}
