package cmd

import (
	"fmt"
	"pegic/executor"
	"pegic/interactive"

	"github.com/desertbit/grumble"
)

func init() {
	usage := `<hashKey>
[ --from <startSortKey> ]
[ --to <stopSortKey> ]
[
  --prefix <filter>
  --suffix <filter>
  --contains <filter>
]`

	scanCmd := &grumble.Command{
		Name:  "scan",
		Help:  "scan records under the hashkey",
		Usage: "\nscan " + usage,
		Run:   requireUseTable(runScanCommand),
		Flags: func(f *grumble.Flags) {
			f.StringL("from", "", "<startSortKey>")
			f.StringL("to", "", "<stopSortKey>")
			f.StringL("prefix", "", "<filter>")
			f.StringL("suffix", "", "<filter>")
			f.StringL("contains", "", "<filter>")
		},
		AllowArgs: true,
	}

	interactive.App.AddCommand(scanCmd)
}

func runScanCommand(c *grumble.Context) error {
	if len(c.Args) < 1 {
		return fmt.Errorf("missing <hashkey> for `scan`")
	}

	from := c.Flags.String("from")
	to := c.Flags.String("to")
	suffix := c.Flags.String("suffix")
	prefix := c.Flags.String("prefix")
	contains := c.Flags.String("contains")

	cmd := &executor.ScanCommand{HashKey: c.Args[0]}
	if from != "" {
		cmd.From = &from
	}
	if to != "" {
		cmd.To = &to
	}
	if suffix != "" {
		cmd.Suffix = &suffix
	}
	if prefix != "" {
		cmd.Prefix = &prefix
	}
	if contains != "" {
		cmd.Contains = &contains
	}
	if err := cmd.Validate(); err != nil {
		return err
	}

	return cmd.IterateAll(globalContext)
}
