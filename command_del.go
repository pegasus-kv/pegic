package pegic

import (
	"context"
	"fmt"
	"pegic/ast"
	p "pegic/parser"
)

type delCommand struct {
	hashKey string
	sortKey string
}

func (c *delCommand) parse(input string) error {
	res, s := p.ArrayWhiteSpace(p.String, p.String)(input)
	if res.Err != nil {
		return res.Err
	}
	if s != "" {
		return fmt.Errorf("redundant input `%s`", s)
	}
	out := res.Output.([]interface{})
	c.hashKey = out[0].(string)
	c.sortKey = out[1].(string)
	return nil
}

func (c *delCommand) execute(ctx *ExecContext) error {
	if ctx.table == nil {
		return noTableError
	}
	if err := ctx.table.Del(context.Background(), []byte(c.hashKey), []byte(c.sortKey)); err != nil {
		return err
	}
	return nil
}

func (*delCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		Arguments: []*ast.CommandArgument{
			{Name: "hashkey"},
			{Name: "sortkey"},
		},
	}
}
