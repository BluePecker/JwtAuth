package main

import (
    "runtime"
    "github.com/BluePecker/JwtAuth/cmd/jwtauthd"
    "fmt"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    jwtauthd.JwtAuth.Cmd.Execute()
    
    fmt.Println(jwtauthd.JwtAuth.Args)
}