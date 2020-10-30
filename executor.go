package pegic

import (
	"fmt"
	"strings"
)

// Executor is the pegic command executor in interactive mode.
func Executor(s string) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return
	}

	parts := strings.SplitN(s, " ", 2)
	cmdStr := strings.ToUpper(parts[0])
	cmd, found := commandsTable[cmdStr]
	if !found {
		fmt.Printf("ERROR: unsupported command: \"%s\"\n", cmdStr)
		return
	}
	var subcommand string = ""
	if len(parts) != 1 {
		subcommand = parts[1]
	}
	err := cmd.parse(strings.TrimSpace(subcommand))
	if err != nil {
		fmt.Printf("ERROR: unable to parse command: %s\n", err)
		return
	}
	if err := cmd.execute(); err != nil {
		fmt.Printf("ERROR: execution failed: %s\n", err)
		return
	}
}
