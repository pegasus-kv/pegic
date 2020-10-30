package pegic

import (
	"fmt"
	"pegic/ast"
	p "pegic/parser"
)

/*
The AST flow. Drawed by asciiflow.com.

         |
         v
      SCAN
         +
         |
    +---------+
    v    |    v
DELETE   | COUNT
    +    |    +
    +---------+
         v
     HASHKEY:arg[1]
         +            +->SUFFIX:arg[1]+------+
         |            |                      |
         +-->SORTKEY--+->PREFIX:arg[1]+------+
         |            |                      |
         |            +->CONTAINS:arg[1]+----+
         |            |                      |
         |            +->BETWEEN:arg[1]      |
         |                  +                |
         |                  +-->AND:arg[1]+--+
         |                                   |
   +-----------------------------------------+
   v     v
 NOVALUE |
   +     |
   +-----+
         v
*/

type scanCommand struct {
	count   bool
	delete  bool
	hashKey string
	sortKey *scanSortKey
	noValue bool
}

type scanSortKey struct {
	start    string
	stop     string
	contains string
	prefix   string
	suffix   string
}

func (c *scanCommand) parse(input string) error {
	res, s := p.ArrayWhiteSpace(
		p.Opt(p.Alt(p.TagNoCase("count"), p.TagNoCase("delete"))),
		p.TagNoCase("hashkey"),
		p.String,
		p.Opt(p.ArrayWhiteSpace(
			p.TagNoCase("sortkey"),
			p.Opt(p.Alt(
				p.ArrayWhiteSpace(p.TagNoCase("between"), p.String, p.TagNoCase("and"), p.String),
				p.ArrayWhiteSpace(p.TagNoCase("contains"), p.String),
				p.ArrayWhiteSpace(p.TagNoCase("prefix"), p.String),
				p.ArrayWhiteSpace(p.TagNoCase("suffix"), p.String),
			)),
		)),
		p.Opt(p.TagNoCase("novalue")))(input)
	if res.Err != nil {
		return res.Err
	}
	if s != "" {
		return fmt.Errorf("redundant input `%s`", s)
	}
	out := res.Output.([]interface{})
	if out[0] != nil {
		if out[0].(string) == "count" {
			c.count = true
			c.delete = false
		} else {
			c.count = false
			c.delete = true
		}
	} else {
		c.count = false
		c.delete = false
	}
	c.hashKey = out[2].(string)
	if out[3] != nil {
		arr := out[3].([]interface{})[1]
		var (
			start    string
			stop     string
			contains string
			prefix   string
			suffix   string
		)
		if arr != nil {
			arr := arr.([]interface{})
			switch arr[0].(string) {
			case "between":
				start = arr[1].(string)
				stop = arr[3].(string)
			case "contains":
				contains = arr[1].(string)
			case "prefix":
				prefix = arr[1].(string)
			case "suffix":
				suffix = arr[1].(string)
			}
		}
		c.sortKey = &scanSortKey{
			start:    start,
			stop:     stop,
			contains: contains,
			prefix:   prefix,
			suffix:   suffix,
		}
	} else {
		c.sortKey = nil
	}
	c.noValue = out[4] != nil
	return nil
}

func (c *scanCommand) execute() error {
	fmt.Printf("%+v\n", c)
	if c.sortKey != nil {
		fmt.Printf("%+v\n", c.sortKey)
	}
	return nil
}

func (*scanCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		SubNodes: map[string]*ast.CommandASTNode{
			"COUNT":  {},
			"DELETE": {},
			"HASHKEY": {
				Arguments: []*ast.CommandArgument{
					{Name: "hashkey"},
				},
			},
			"SORTKEY": {
				SubNodes: map[string]*ast.CommandASTNode{
					"SUFFIX": {
						Arguments: []*ast.CommandArgument{
							{Name: "suffix"},
						},
					},
					"PREFIX": {
						Arguments: []*ast.CommandArgument{
							{Name: "prefix"},
						},
					},
					"CONTAINS": {
						Arguments: []*ast.CommandArgument{
							{Name: "contains"},
						},
					},
					"BETWEEN": {
						Arguments: []*ast.CommandArgument{
							{Name: "start"},
						},
						SubNodes: map[string]*ast.CommandASTNode{
							"AND": {
								Arguments: []*ast.CommandArgument{
									{Name: "end"},
								},
							},
						},
					},
				},
			},
			"NOVALUE": {},
		},
	}
}
