package main

import (
    "runtime"
    "github.com/BluePecker/JwtAuth/cmd/jwtauthd"
    "github.com/BluePecker/JwtAuth/daemon"
    "fmt"
    //"github.com/BluePecker/JwtAuth/server"
    //"github.com/BluePecker/JwtAuth/server/router/jwt"
)

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    if err := jwtauthd.JwtAuth.Cmd.Execute(); err == nil {
        
        fmt.Println(jwtauthd.JwtAuth.Args)
        
        (&daemon.Daemon{}).Start(jwtauthd.JwtAuth.Args)
        
        //api := &server.Server{}
        //api.AddRouter(jwt.NewRouter(nil))
        //
        //api.Accept(server.Options{Host: "", Port: jwtauthd.JwtAuth.Args.Port})
    } else {
        fmt.Println(err)
    }
}