package pegic

import (
	"context"
	"fmt"
	"os"
	"pegic/ast"
	"pegic/parser"
)

// pegicCommand is the interface for all the supported commands in pegic.
type pegicCommand interface {
	// Parse the input string after first keyword
	parse(input string) error

	// Execute the command
	execute(*ExecContext) error

	// Create a AST node for this command.
	astNode() *ast.CommandASTNode
}

var commandsTable = map[string]pegicCommand{
	"USE":      &useCommand{},
	"LS":       &lsCommand{},
	"GET":      &getCommand{},
	"SET":      &setCommand{},
	"DEL":      &delCommand{},
	"SCAN":     &scanCommand{},
	"FULLSCAN": &fullScanCommand{},
	// "COMPRESSION": &compressionCommand{},
	// "ENCODING":    &encodingCommand{},
	// "EXIT":        &exitCommand{},
}

type useCommand struct {
	tableName string
}

func (c *useCommand) parse(input string) error {
	res, s := parser.String(input)
	if res.Err != nil {
		return res.Err
	}
	if s != "" {
		return fmt.Errorf("redundant input `%s`", s)
	}
	c.tableName = res.Output.(string)
	return nil
}

func (c *useCommand) execute(ctx *ExecContext) error {
	tb, err := ctx.client.OpenTable(context.Background(), c.tableName)
	if err != nil {
		return err
	}
	ctx.table = tb
	return nil
}

func (*useCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		Arguments: []*ast.CommandArgument{
			{Name: "table_name"},
		},
	}
}

type exitCommand struct {
}

func (c *exitCommand) parse(input string) error {
	return nil
}

func (*exitCommand) execute(*ExecContext) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}

func (*exitCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		CustomDescription: "Exit this program",
	}
}
