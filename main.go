package main

import (
    "runtime"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/action"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    if err := action.RootCmd.Cmd.Execute(); err == nil {
        // todo
    } else {
        logrus.Error(err)
    }
}