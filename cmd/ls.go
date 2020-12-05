package cmd

import (
	"pegic/interactive"

	"github.com/desertbit/grumble"
)

func init() {
	interactive.App.AddCommand(&grumble.Command{
		Name:    "ls",
		Aliases: []string{"LS"},
		Help:    "list the tables in the cluster",
		Run: func(c *grumble.Context) error {
			return nil
		},
		AllowArgs: true,
	})
}
