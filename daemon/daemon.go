package daemon

import (
    "os"
    "fmt"
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/pkg/storage"
    "github.com/BluePecker/JwtAuth/daemon/server"
    "syscall"
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
    
    StorageE *storage.Engine
}

func Logger(level string) {
    logrus.SetFormatter(&logrus.TextFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })
    Level, err := logrus.ParseLevel(level)
    if err != nil {
        logrus.Error(err)
        os.Exit(0)
    }
    logrus.SetLevel(Level)
}

func Version(version bool) {
    if version == true {
        fmt.Printf("JwtAuth version %s.\n", VERSION)
        os.Exit(0)
    }
}

func NewDaemon(background bool, args Options) (*Daemon, *daemon.Context) {
    
    if background {
        ctx := &daemon.Context{
            PidFileName: args.PidFile,
            PidFilePerm: 0644,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
            LogFileName: args.LogFile,
        }
        if process, err := ctx.Reborn(); err == nil {
            if process != nil {
                return nil, nil
            } else {
                return &Daemon{Options: &args}, ctx
            }
        } else {
            if err == daemon.ErrWouldBlock {
                logrus.Error("daemon already exists.")
            } else {
                logrus.Errorf("Unable to run: ", err)
            }
            os.Exit(0)
        }
    }
    return &Daemon{Options: &args}, nil
}

func NewStart(args Options) {
    
    Logger(args.LogLevel)
    Version(args.Version)
    
    if args.Secret == "" {
        fmt.Println("please specify secret for jwt encode.")
        os.Exit(0)
    }
    
    if progress, ctx := NewDaemon(args.Daemon, args); progress == nil {
        return
    } else {
        if (ctx != nil) {
            defer ctx.Release()
        }
        if err := progress.Storage(); err != nil {
            logrus.Error(err)
            os.Exit(0)
        }
        
        quit := make(chan struct{})
        go progress.Shadow(quit)
        go func() {
            go progress.Rosiness(quit)
            defer logrus.Infof("now listening on: http://%s:%d", args.Host, args.Port)
        }()
        daemon.SetSigHandler(func(sig os.Signal) error {
            close(quit)
            return daemon.ErrStop
        }, syscall.SIGTERM, syscall.SIGQUIT)
        
        if err := daemon.ServeSignals(); err != nil {
            logrus.Error(err)
        }
        logrus.Error("jwt daemon terminated.")
    }
}