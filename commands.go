package pegic

import (
	"fmt"
	"os"
	"pegic/ast"
)

// pegicCommand is the interface for all the supported commands in pegic.
type pegicCommand interface {

	// Execute the command
	execute(parsedCmd *ast.ParsedCommand) error

	// Create a AST node for this command.
	astNode() *ast.CommandASTNode
}

var commandsTable = map[string]pegicCommand{
	"USE":         &useCommand{},
	"SET":         &setCommand{},
	"COMPRESSION": &compressionCommand{},
	"ENCODING":    &encodingCommand{},
	"EXIT":        &exitCommand{},
}

type useCommand struct {
	tableName string
}

func (*useCommand) execute(parsedCmd *ast.ParsedCommand) error {
	// TODO
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

func (*exitCommand) execute(parsedCmd *ast.ParsedCommand) error {
	fmt.Println("Bye!")
	os.Exit(0)
	return nil
}

func (*exitCommand) astNode() *ast.CommandASTNode {
	return &ast.CommandASTNode{
		CustomDescription: "Exit this program",
	}
}
