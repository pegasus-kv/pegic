package cmd

import (
	"pegic/interactive"

	"github.com/desertbit/grumble"
)

func init() {
	interactive.App.AddCommand(&grumble.Command{
		Name:    "get",
		Aliases: []string{"GET"},
		Help:    "read a record from Pegasus",
		Usage:   "get <HASHKEY> <SORTKEY>",
		Run: func(c *grumble.Context) error {
			// TODO(wutao): verify if the use table exists
			c.App.Println("ok")
			return nil
		},
		AllowArgs: true,
	})
}
