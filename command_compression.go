package pegic

import (
	"fmt"
	"pegic/ast"
	p "pegic/parser"
)

type compressionCommand struct {
	algorithm string
}

func (c *compressionCommand) parse(input string) error {
	res, s := p.Alt(p.TagNoCase("zstd"), p.TagNoCase("no"))(input)
	if res.Err != nil {
		return res.Err
	}
	if s != "" {
		return fmt.Errorf("redundant input `%s`", s)
	}
	c.algorithm = res.Output.(string)
	return nil
}

func (*compressionCommand) execute() error {
	// TODO(wutao)
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
