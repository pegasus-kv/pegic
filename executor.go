package pegic

import (
	"fmt"
	"pegic/ast"
	"strings"
)

// Executor is the pegic command executor in interactive mode.
func Executor(s string) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return
	}

	args := strings.Split(s, " ")
	cmdStr := strings.ToUpper(args[0])
	cmd, found := commandsTable[cmdStr]
	if !found {
		fmt.Printf("ERROR: unsupported command: \"%s\"\n", cmdStr)
		return
	}
	parsedCmd, err := ast.Parse(args)
	if err != nil {
		fmt.Printf("ERROR: unable to parse command: %s\n", err)
		return
	}
	if err := cmd.execute(parsedCmd); err != nil {
		fmt.Printf("ERROR: execution failed: %s\n", err)
		return
	}
}
