package interactive

import (
	"github.com/desertbit/grumble"
)

// App is the global interactive application.
var App *grumble.App

// TODO(jiashuo1) some command with table no support verify the table exists
func init() {
	App = grumble.New(&grumble.Config{
		Name: "pegic",
	})
}

func Run() {
	App.Run()
}
