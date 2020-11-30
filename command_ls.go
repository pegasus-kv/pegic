package pegic

import (
	"context"
	"fmt"
	"pegic/ast"
)

type lsCommand struct {
}

func (c *lsCommand) parse(string) error {
	return nil
}

func (c *lsCommand) execute(ctx *ExecContext) error {
	tables, err := ctx.adminClient.ListTables(context.Background())
	if err != nil {
		return err
	}
	for _, table := range tables {
		fmt.Println(table.Name)
	}
	return nil
}

func (c *lsCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		CustomDescription: "List all tables",
	}
}
