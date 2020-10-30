package pegic

import (
	"fmt"
	"pegic/ast"
	p "pegic/parser"
)

type fullScanCommand struct {
	count   bool
	delete  bool
	hashKey *fullScanHashKey
	sortKey *fullScanSortKey
	noValue bool
}

type fullScanHashKey struct {
	contains string
	suffix   string
	prefix   string
}

type fullScanSortKey struct {
	is       string
	contains string
	suffix   string
	prefix   string
	start    string
	stop     string
}

func (c *fullScanCommand) parse(input string) error {
	res, s := p.ArrayWhiteSpace(
		p.Opt(p.Alt(p.TagNoCase("count"), p.TagNoCase("delete"))),
		p.Opt(
			p.ArrayWhiteSpace(
				p.TagNoCase("hashkey"),
				p.Alt(p.TagNoCase("contains"), p.TagNoCase("suffix"), p.TagNoCase("prefix")),
				p.String,
			),
		),
		p.Opt(
			p.ArrayWhiteSpace(
				p.TagNoCase("sortkey"),
				p.Opt(p.Alt(
					p.ArrayWhiteSpace(p.Alt(p.TagNoCase("is"), p.TagNoCase("contains"), p.TagNoCase("suffix"), p.TagNoCase("prefix")), p.String),
					p.ArrayWhiteSpace(p.TagNoCase("between"), p.String, p.TagNoCase("and"), p.String),
				)),
			),
		),
		p.Opt(p.TagNoCase("novalue")),
	)(input)
	if res.Err != nil {
		return res.Err
	}
	if s != "" {
		return fmt.Errorf("redundant input `%s`", s)
	}
	out := res.Output.([]interface{})
	if out[0] != nil {
		c.count = out[0].(string) == "count"
		c.delete = !c.count
	} else {
		c.count = false
		c.delete = false
	}
	if out[1] != nil {
		hk := &fullScanHashKey{}
		p := out[1].([]interface{})
		switch p[1].(string) {
		case "contains":
			hk.contains = p[2].(string)
		case "suffix":
			hk.suffix = p[2].(string)
		case "prefix":
			hk.prefix = p[2].(string)
		}
		c.hashKey = hk
	} else {
		c.hashKey = nil
	}
	if out[2] != nil {
		sk := &fullScanSortKey{}
		p := out[2].([]interface{})
		if p[1] != nil {
			a := p[1].([]interface{})
			if a[0] == "between" {
				sk.start = a[1].(string)
				sk.stop = a[3].(string)
			} else {
				switch a[0].(string) {
				case "is":
					sk.is = a[1].(string)
				case "contains":
					sk.contains = a[1].(string)
				case "prefix":
					sk.prefix = a[1].(string)
				case "suffix":
					sk.suffix = a[1].(string)
				}
			}
		}
		c.sortKey = sk
	} else {
		c.sortKey = nil
	}
	c.noValue = out[3] != nil
	return nil
}

func (c *fullScanCommand) execute() error {
	fmt.Printf("%+v\n", c)
	if c.hashKey != nil {
		fmt.Printf("%+v\n", c.hashKey)
	}
	if c.sortKey != nil {
		fmt.Printf("%+v\n", c.sortKey)
	}
	return nil
}

func (c *fullScanCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		SubNodes: map[string]*ast.CommandASTNode{
		},
	}
}
