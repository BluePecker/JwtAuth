package main

import (
    "runtime"
    "github.com/BluePecker/JwtAuth/cmd/jwtauthd"
    "fmt"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    if err := jwtauthd.JwtAuth.Cmd.Execute(); err == nil {
    
        fmt.Println(jwtauthd.JwtAuth.Args)
    } else {
        fmt.Println(err)
    }
}