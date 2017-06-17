package daemon

import (
    "github.com/BluePecker/JwtAuth/cmd/jwtauthd"
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/server"
    "github.com/BluePecker/JwtAuth/server/router/jwt"
    "fmt"
)

type Daemon struct {
    
}

func (d *Daemon) Start(conf jwtauthd.Args) {
    fmt.Println(conf.Daemon, conf.PidFile, conf.LogFile)
    if (conf.Daemon == true) {
        dCtx := daemon.Context{
            PidFileName: conf.PidFile,
            PidFilePerm: 0644,
            LogFileName: conf.LogFile,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
        }
        
        if child, err := dCtx.Reborn(); err != nil {
            logrus.Fatal(err)
        } else if child != nil {
            return
        }
        
        defer dCtx.Release()
    }
    
    api := &server.Server{}
    api.AddRouter(jwt.NewRouter(nil))
    
    api.Accept(server.Options{Host: "", Port: conf.Port})
}