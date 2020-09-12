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

	args := strings.Split(s, " ")
	cmdStr := strings.ToUpper(args[0])
	cmd, found := commandsTable[cmdStr]
	if !found {
		fmt.Printf("unsupported command: \"%s\"\n", cmdStr)
		return
	}
	if err := cmd.execute(); err != nil {
		fmt.Printf("execution failed: %s\n", err)
	}
}
