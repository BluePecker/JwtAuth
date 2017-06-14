package main

import (
    "runtime"
    "github.com/BluePecker/JwtAuth/cmd/jwtauthd"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    jwtauthd.JwtAuth.Cmd.Execute()
}