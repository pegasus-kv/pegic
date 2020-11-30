package pegic

import (
	"context"
	"fmt"
	"pegic/ast"
	p "pegic/parser"
)

type getCommand struct {
	hashKey string
	sortKey string
}

func (c *getCommand) parse(input string) error {
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

func (c *getCommand) execute(ctx *ExecContext) error {
	if ctx.table == nil {
		return noTableError
	}
	res, err := ctx.table.Get(context.Background(), []byte(c.hashKey), []byte(c.sortKey))
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

func (c *getCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		Arguments: []*ast.CommandArgument{
			{Name: "hashkey"},
			{Name: "sortkey"},
		},
	}
}
