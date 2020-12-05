package main

import (
	"pegic/cmd"
	"pegic/shell"

	"github.com/XiaoMi/pegasus-go-client/pegalog"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// pegasus-go-client's logs use the same logger as ours.
	pegalog.SetLogger(logrus.StandardLogger())
	// configure log destination
	logrus.SetOutput(&lumberjack.Logger{
		Filename:  "./pegic.log",
		LocalTime: true,
	})

	cmd.Init()
	shell.Root.Execute()
}