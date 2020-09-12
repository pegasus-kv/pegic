package pegic

import (
	"pegic/ast"

	"github.com/c-bata/go-prompt"
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
			{Name: "zstd|no", Suggester: func() []prompt.Suggest {
				return []prompt.Suggest{
					{Text: "zstd"},
					{Text: "no"},
				}
			}},
		},
	}
}
