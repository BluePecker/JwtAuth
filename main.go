package main

import (
    "runtime"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/cmd/jwtauthd"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    if err := jwtauthd.RootCmd.Cmd.Execute(); err == nil {
        // todo
    } else {
        logrus.Error(err)
    }
}