package ast

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

// CommandArgument is a required field followed after a KEYWORD.
type CommandArgument struct {
	Name      string
	Suggester func() []prompt.Suggest

	value string
}

// CommandASTNode represents a node in the AST for interactive command parsing.
type CommandASTNode struct {
	Arguments []*CommandArgument

	// If no CustomDescription is given, the description shows all the arguments followed.
	CustomDescription string

	SubNodes map[string]*CommandASTNode
}

func (n *CommandASTNode) description() string {
	if n.CustomDescription != "" {
		return n.CustomDescription
	}
	dest := ""
	for _, arg := range n.Arguments {
		dest += fmt.Sprintf("<%s> ", arg.Name)
	}
	return dest
}

var root *CommandASTNode = &CommandASTNode{
	SubNodes: make(map[string]*CommandASTNode),
}

// RegisterCommand registers a CommandASTNode to the AST tree.
func RegisterCommand(cmdName string, cmdNode *CommandASTNode) {
	root.SubNodes[cmdName] = cmdNode
}

// Suggest returns the suggestions given with the user input.
func Suggest(arguments []string) []prompt.Suggest {
	return suggestInTree(arguments, root)
}

func suggestInTree(arguments []string, parent *CommandASTNode) []prompt.Suggest {
	word := arguments[0]
	if len(arguments) == 1 { // user is typing the keyword
		var sgts []prompt.Suggest
		for keyword, subNode := range parent.SubNodes {
			sgts = append(sgts, prompt.Suggest{
				Text:        keyword,
				Description: subNode.description(),
			})
		}
		return prompt.FilterHasPrefix(sgts, word, true)
	}

	node, found := parent.SubNodes[word]
	if !found { // the keyword is invalid
		return []prompt.Suggest{}
	}

	arguments = arguments[1:]
	if len(node.Arguments) > 0 { // the keyword is followed after some arguments
		if len(arguments) <= len(node.Arguments) {
			// user is typing the argument
			cmdArg := node.Arguments[len(arguments)-1]
			if cmdArg.Suggester != nil {
				return cmdArg.Suggester()
			}
			return []prompt.Suggest{}
		}
		arguments = arguments[:len(arguments)-len(node.Arguments)]
	}
	if len(arguments) > 0 {
		return suggestInTree(arguments, node)
	}
	return []prompt.Suggest{}
}
