package daemon

import (
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/server"
    "github.com/BluePecker/JwtAuth/server/router/jwt"
)

var (
    // Redis/Mongodb connection pool size
    MaxPoolSize int = 50
)

type Storage struct {
    Driver   string
    Path     string
    Host     string
    Port     int
    Username string
    Password string
    PoolSize int
}

type Security struct {
    TLS  bool
    Key  string
    Cert string
}

type Options struct {
    PidFile string
    LogFile string
    
    Port    int
    Host    string
    
    Daemon  bool
    
    Https   Security
}

// todo
type Daemon struct{}

func NewStart(options Options) {
    if (options.Daemon == true) {
        dCtx := daemon.Context{
            PidFileName: options.PidFile,
            PidFilePerm: 0644,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
            LogFileName: options.LogFile,
        }
        
        defer dCtx.Release()
        
        if child, err := dCtx.Reborn(); err != nil {
            logrus.Error(err)
        } else if child != nil {
            return
        }
    }
    
    api := &server.Server{}
    api.AddRouter(jwt.NewRouter(nil))
    
    var TLS *server.TLS
    if options.Https.TLS {
        TLS = &server.TLS{
            Key: options.Https.Key,
            Cert: options.Https.Cert,
        }
    }
    
    api.Accept(server.Options{Host: options.Host, Tls: TLS, Port: options.Port})
}