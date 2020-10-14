package pegic

import "pegic/ast"

type setCommand struct {
	hashKey    string
	sortKey    string
	value      string
	ttlSeconds string
}

func (*setCommand) execute(parsedCmd *ast.ParsedCommand) error {
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
