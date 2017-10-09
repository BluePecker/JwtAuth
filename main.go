package main

import (
	"runtime"
	"github.com/Sirupsen/logrus"
	"github.com/BluePecker/JwtAuth/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := cmd.RootCmd.Cmd.Execute(); err == nil {
		// todo
	} else {
		logrus.Error(err)
	}
}
