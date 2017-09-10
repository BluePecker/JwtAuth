package daemon

import (
    "os"
    "fmt"
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/pkg/storage"
    "github.com/BluePecker/JwtAuth/daemon/server"
)

const (
    TOKEN_TTL = 2 * 3600
    
    VERSION = "1.0.0"
    
    ALLOW_LOGIN_NUM = 3
)

type Storage struct {
    Driver string
    Opts   string
}

type TLS struct {
    Key  string
    Cert string
}

type Options struct {
    PidFile  string
    LogFile  string
    LogLevel string
    Port     int
    Host     string
    SockFile string
    Daemon   bool
    Version  bool
    TLS      TLS
    Storage  Storage
    Secret   string
}

type Daemon struct {
    Options  *Options
    
    shadow   *server.Shadow
    rosiness *server.Rosiness
    
    Store    *storage.Driver
}

func NewDaemon(background bool, args Options) *Daemon {
    if background {
        ctx := daemon.Context{
            PidFileName: args.PidFile,
            PidFilePerm: 0644,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
            LogFileName: args.LogFile,
        }
        if rank, err := logrus.ParseLevel(args.LogLevel); err != nil {
            fmt.Println(err)
            os.Exit(0)
        } else {
            logrus.SetFormatter(&logrus.TextFormatter{
                TimestampFormat: "2006-01-02 15:04:05",
            })
            logrus.SetLevel(rank)
        }
        if process, err := ctx.Reborn(); err == nil {
            defer ctx.Release()
            if process != nil {
                return nil
            }
        } else {
            if err == daemon.ErrWouldBlock {
                fmt.Println("daemon already exists.")
            } else {
                fmt.Println("Unable to run: ", err)
            }
            os.Exit(0)
        }
    }
    return &Daemon{Options: &args}
}

func NewStart(args Options) {
    if args.Version == true {
        fmt.Printf("JwtAuth version %s.\n", VERSION)
        os.Exit(0)
    }
    
    if args.Secret == "" {
        fmt.Println("please specify the key.")
        os.Exit(0)
    }
    proc := NewDaemon(args.Daemon, args)
    if proc == nil {
        return
    }
    if err := proc.Storage(); err != nil {
        logrus.Error(err)
        os.Exit(0)
    }
    go func() {
        proc.Shadow()
    }()
    if err := proc.Rosiness(); err != nil {
        logrus.Error(err)
    } else {
        logrus.Error("api server closed.");
    }
    os.Exit(0)
}