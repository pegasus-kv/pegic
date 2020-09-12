package pegic

import (
	"pegic/ast"
	"strings"

	"github.com/c-bata/go-prompt"
)

func init() {
	for cmdName, cmd := range commandsTable {
		ast.RegisterCommand(cmdName, cmd.astNode())
	}
}

// Completer is the pegic auto-completer in interactive mode.
func Completer(d prompt.Document) []prompt.Suggest {
	args := strings.Split(d.TextBeforeCursor(), " ")
	return ast.Suggest(args)
}
