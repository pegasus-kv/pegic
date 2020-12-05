package cmd

import (
	"pegic/interactive"

	"github.com/desertbit/grumble"
)

func init() {
	interactive.App.AddCommand(&grumble.Command{
		Name:    "set",
		Aliases: []string{"SET"},
		Help:    "write a record into Pegasus",
		Usage:   "set <HASHKEY> <SORTKEY> <VALUE>",
		Run: func(c *grumble.Context) error {
			// TODO(wutao): verify if the use table exists
			c.App.Println("ok")
			return nil
		},
		AllowArgs: true,
	})
}
