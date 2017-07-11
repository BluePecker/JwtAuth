package daemon

import (
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/server"
    "github.com/BluePecker/JwtAuth/server/router/jwt"
)

type Conf struct {
    PidFile string
    LogFile string
    
    Port    int
    Host    string
    
    TLS     bool
    TLSKey  string
    TLSCert string
    
    Daemon  bool
}

type Daemon struct {
    
}

func NewStart(conf Conf) {
    if (conf.Daemon == true) {
        dCtx := daemon.Context{
            PidFileName: conf.PidFile,
            PidFilePerm: 0644,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
            LogFileName: conf.LogFile,
        }
        
        defer dCtx.Release()
        
        if child, err := dCtx.Reborn(); err != nil {
            logrus.Fatal(err)
        } else if child != nil {
            return
        }
    }
    
    api := &server.Server{}
    api.AddRouter(jwt.NewRouter(nil))
    
    var TLS *server.TLS
    if conf.TLS {
        TLS = &server.TLS{
            Cert: conf.TLSCert,
            Key: conf.TLSKey,
        }
    }
    
    api.Accept(server.Options{Host: conf.Host, Tls: TLS, Port: conf.Port})
}