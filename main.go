package main

import (
    "runtime"
    "github.com/BluePecker/JwtAuth/cmd/jwtauthd"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    if err := jwtauthd.JwtAuth.Cmd.Execute(); err != nil {
        panic(err)
    }
}