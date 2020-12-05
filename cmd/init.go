package cmd

import (
	"fmt"
	"os"
	"pegic/executor"

	"github.com/desertbit/grumble"
)

var globalContext *executor.Context

func Init(metaAddrs []string) {
	globalContext = executor.NewContext(os.Stdout, metaAddrs)
}

func requireUseTable(run func(*grumble.Context) error) func(c *grumble.Context) error {
	grumbleRun := func(c *grumble.Context) error {
		if globalContext.UseTable == nil {
			return fmt.Errorf("please USE a table first")
		}
		return run(c)
	}
	return grumbleRun
}
