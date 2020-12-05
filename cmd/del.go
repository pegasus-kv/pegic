package cmd

import (
	"pegic/interactive"

	"github.com/desertbit/grumble"
)

func init() {
	interactive.App.AddCommand(&grumble.Command{
		Name:    "del",
		Aliases: []string{"DEL"},
		Help:    "delete a record",
		Usage:   "del <HASHKEY> <SORTKEY>",
		Run: func(c *grumble.Context) error {
			return nil
		},
		AllowArgs: true,
	})
}
