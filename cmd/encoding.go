package cmd

import (
	"pegic/interactive"

	"github.com/desertbit/grumble"
)

func init() {
	rootCmd := &grumble.Command{
		Name:    "encoding",
		Aliases: []string{"ENCODING"},
		Help:    "read the current encoding",
		Run: func(c *grumble.Context) error {
			// TODO(wutao): verify if the use table exists
			c.App.Println("ok")
			return nil
		},
		AllowArgs: true,
	}

	rootCmd.AddCommand(&grumble.Command{
		Name:    "hashkey",
		Aliases: []string{"HASHKEY"},
		Help:    "set encoding for hashkey",
		Run: func(c *grumble.Context) error {
			// TODO(wutao): verify if the use table exists
			c.App.Println("ok")
			return nil
		},
		AllowArgs: true,
	})

	rootCmd.AddCommand(&grumble.Command{
		Name:    "sortkey",
		Aliases: []string{"SORTKEY"},
		Help:    "set encoding for sortkey",
		Run: func(c *grumble.Context) error {
			// TODO(wutao): verify if the use table exists
			c.App.Println("ok")
			return nil
		},
		AllowArgs: true,
	})

	rootCmd.AddCommand(&grumble.Command{
		Name:    "value",
		Aliases: []string{"VALUE"},
		Help:    "set encoding for value",
		Run: func(c *grumble.Context) error {
			// TODO(wutao): verify if the use table exists
			c.App.Println("ok")
			return nil
		},
		AllowArgs: true,
	})

	interactive.App.AddCommand(rootCmd)
}
