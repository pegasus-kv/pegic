package pegic

import (
	"pegic/ast"
)

type compressionCommand struct {
	algorithm string
}

func (*compressionCommand) execute() error {
	return nil
}

func (*compressionCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		Arguments: []*ast.CommandArgument{
			{Name: "zstd|no", Selections: &ast.Select{
				Items: []string{
					"zstd",
					"no",
				},
			}},
		},
	}
}
