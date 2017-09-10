package main

import (
    "runtime"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/cmd/jwtd"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    if err := jwtd.RootCmd.Cmd.Execute(); err == nil {
        // todo
    } else {
        logrus.Error(err)
    }
}