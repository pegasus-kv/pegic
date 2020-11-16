package pegic

import (
	"context"
	"fmt"
	"pegic/ast"
	p "pegic/parser"
	"time"
)

type setCommand struct {
	hashKey    string
	sortKey    string
	value      string
	ttlSeconds uint
}

func (c *setCommand) parse(input string) error {
	res, s := p.ArrayWhiteSpace(p.String, p.String, p.String, p.Opt(p.UInt))(input)
	if res.Err != nil {
		return res.Err
	}
	if s != "" {
		return fmt.Errorf("redundant input `%s`", s)
	}
	out := res.Output.([]interface{})
	c.hashKey = out[0].(string)
	c.sortKey = out[1].(string)
	c.value = out[2].(string)
	if out[3] != nil {
		c.ttlSeconds = out[3].(uint)
	} else {
		c.ttlSeconds = 0
	}
	return nil
}

func (c *setCommand) execute(ctx *ExecContext) error {
	if ctx.table == nil {
		return noTableError
	}
	if c.ttlSeconds != 0 {
		if err := ctx.table.SetTTL(context.Background(), []byte(c.hashKey), []byte(c.sortKey), []byte(c.value), time.Duration(c.ttlSeconds) * time.Second); err != nil {
			return err
		}
	} else {
		if err := ctx.table.Set(context.Background(), []byte(c.hashKey), []byte(c.sortKey), []byte(c.value)); err != nil {
			return err
		}
	}
	return nil
}

func (*setCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		Arguments: []*ast.CommandArgument{
			{Name: "hashkey"},
			{Name: "sortkey"},
			{Name: "value"},
			{Name: "ttlSeconds"},
		},
	}
}
