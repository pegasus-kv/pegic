package shell

import (
	"pegic/interactive"

	"github.com/spf13/cobra"
)

// Root command for pegic.
var Root *cobra.Command

func init() {
	Root = &cobra.Command{
		Use:   "pegic [--meta|-m <meta-list>]",
		Short: "pegic: Pegasus Interactive Command-Line tool",
		PreRun: func(cmd *cobra.Command, args []string) {
			// validate meta-list
		},
		Run: func(cmd *cobra.Command, args []string) {
			// the default entrance is interactive
			interactive.Run()
		},
	}
}
