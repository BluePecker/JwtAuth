package daemon

import (
    "reflect"
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/server"
    "github.com/BluePecker/JwtAuth/storage"
    "github.com/BluePecker/JwtAuth/server/router/jwt"
    "github.com/BluePecker/JwtAuth/server/router"
    "os"
    _ "github.com/BluePecker/JwtAuth/storage/redis"
    _ "github.com/BluePecker/JwtAuth/storage/ram"
    "fmt"
)

const (
    VERSION = "1.0.0"
)

type Storage struct {
    Driver     string
    Path       string
    Host       string
    Port       int
    MaxRetries int
    Username   string
    Password   string
    PoolSize   int
    Database   string
}

type Security struct {
    TLS  bool
    Key  string
    Cert string
}

type Options struct {
    PidFile  string
    LogFile  string
    LogLevel string
    Port     int
    Host     string
    Daemon   bool
    Version  bool
    Security Security
    Storage  Storage
}

type Daemon struct {
    Options *Options
    Server  *server.Server
    Storage *storage.Driver
}

func (d *Daemon) storageOptionInject(p2 *storage.Option) {
    p1 := &(d.Options.Storage)
    u1 := reflect.ValueOf(p1).Elem()
    u2 := reflect.ValueOf(p2).Elem()
    
    for seq := 0; seq < u2.NumField(); seq++ {
        item := u2.Type().Field(seq)
        v1 := u1.FieldByName(item.Name)
        v2 := u2.FieldByName(item.Name)
        if v1.IsValid() {
            if v2.Type() == v1.Type() {
                v2.Set(v1)
            }
        }
    }
}

func (d *Daemon) NewStorage() (*storage.Driver, error) {
    option := &storage.Option{}
    d.storageOptionInject(option)
    driver, err := storage.New(d.Options.Storage.Driver, *option)
    return &driver, err
}

func (d *Daemon) NewServer() {
    d.Server = &server.Server{}
}

func (d *Daemon) Listen() {
    if d.Server == nil {
        d.NewServer()
    }
    
    options := server.Options{
        Host: d.Options.Host,
        Port: d.Options.Port,
    }
    
    if d.Options.Security.TLS {
        options.Tls = &server.TLS{
            Cert: d.Options.Security.Cert,
            Key: d.Options.Security.Key,
        }
    }
    
    d.Server.Accept(options)
}

func (d *Daemon) addRouter(routers... router.Router) {
    if d.Server == nil {
        d.NewServer()
    }
    for _, route := range routers {
        d.Server.AddRouter(route)
    }
}

func NewStart(args Options) {
    var err error;
    
    if args.Version == true {
        fmt.Printf("JwtAuth version %s.\n", VERSION)
        os.Exit(0)
    }
    
    if args.Daemon == true {
        dCtx := daemon.Context{
            PidFileName: args.PidFile,
            PidFilePerm: 0644,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
            LogFileName: args.LogFile,
        }
        
        level, err := logrus.ParseLevel(args.LogLevel)
        if err == nil {
            logrus.SetLevel(level)
            logrus.SetFormatter(&logrus.TextFormatter{
                TimestampFormat: "2006-01-02 15:04:05",
            })
        } else {
            logrus.Fatal(err)
        }
        defer dCtx.Release()
        
        if child, err := dCtx.Reborn(); err != nil {
            logrus.Fatal(err)
        } else if child != nil {
            return
        }
    }
    
    jwtPro := &Daemon{
        Options: &args,
    }
    
    if jwtPro.Storage, err = jwtPro.NewStorage(); err != nil {
        logrus.Error(err)
        os.Exit(0)
    }
    
    jwtPro.addRouter(jwt.NewRouter(nil))
    jwtPro.Listen()
}