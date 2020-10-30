package pegic

import (
	"fmt"
	"pegic/ast"
	p "pegic/parser"
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

func (c *setCommand) execute() error {
	fmt.Printf("%+v\n", c)
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
