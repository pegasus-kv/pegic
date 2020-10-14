package ast

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

//
type Select struct {
	Items []string
}

//
func (s *Select) Suggest() []prompt.Suggest {
	res := []prompt.Suggest{}
	for _, item := range s.Items {
		res = append(res, prompt.Suggest{
			Text: item,
		})
	}
	return res
}

// Contains returns whether inputArg matches one of the selection items.
func (s *Select) Contains(inputArg string) bool {
	for _, item := range s.Items {
		if inputArg == item {
			return true
		}
	}
	return false
}

// CommandArgument is a required field followed after a KEYWORD.
type CommandArgument struct {
	Name string

	// Optional.
	Selections *Select
}

// CommandASTNode represents a node in the AST for interactive command parsing.
type CommandASTNode struct {
	// For example:
	// HASHKEY <UTF-8|INT|ASCII>
	// HASHKEY is the keyword, <...> reprensents an argument having a list of options.
	// In this case, len(Arguments)==1.
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

// Root node of the AST.
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
			// suggest with all possible keywords
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
			if cmdArg.Selections != nil {
				return cmdArg.Selections.Suggest()
			}
			// no suggestion for command argument
			return []prompt.Suggest{}
		}
		arguments = arguments[:len(arguments)-len(node.Arguments)]
	}
	if len(arguments) > 0 { // step forward to the next keyword
		return suggestInTree(arguments, node)
	}
	// Keyword is typed completely, no other arguments required.
	// No prompt.
	return []prompt.Suggest{}
}

// ParsedCommand is the parsed result for user input.
type ParsedCommand struct {
	Tokens map[string]([]*ParsedArgument)
}

//
type ParsedArgument struct {
	Value string
}

// Parse out user input into a complete command.
func Parse(arguments []string) (*ParsedCommand, error) {
	parsedCmd := &ParsedCommand{
		Tokens: make(map[string]([]*ParsedArgument)),
	}

	node := root
	for len(arguments) > 0 {
		keyword := arguments[0]
		subNode, hasNode := node.SubNodes[keyword]
		if !hasNode {
			return nil, fmt.Errorf("\"%s\" is not a valid keyword", keyword)
		}
		if len(subNode.Arguments) > len(arguments) {
			return nil, fmt.Errorf("not enough arguments after keyword \"%s\"", keyword)
		}
		token := parsedCmd.Tokens[keyword]
		for i, arg := range subNode.Arguments {
			inputArg := arguments[i]
			// validate argument
			if arg.Selections != nil && !arg.Selections.Contains(inputArg) {
				return nil, fmt.Errorf("argument \"%s\" after keyword \"%s\" must be one of %s", inputArg, keyword, arg.Selections.Items)
			}
			token = append(token, &ParsedArgument{Value: inputArg})
		}

		// step down to the next level
		arguments = arguments[len(subNode.Arguments):]
		node = subNode
	}

	return parsedCmd, nil
}
