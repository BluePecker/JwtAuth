package daemon

import (
    "os"
    "time"
    "fmt"
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/server/types/token"
    "github.com/BluePecker/JwtAuth/server"
    "github.com/BluePecker/JwtAuth/storage"
    //"github.com/BluePecker/JwtAuth/server/router"
    _ "github.com/BluePecker/JwtAuth/storage/redis"
    //_ "github.com/BluePecker/JwtAuth/storage/ram"
    "github.com/dgrijalva/jwt-go"
    RouteToken "github.com/BluePecker/JwtAuth/server/router/token"
    "github.com/kataras/iris"
    "github.com/kataras/iris/core/netutil"
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

type Security struct {
    Verify bool
    TLS    bool
    Key    string
    Cert   string
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
    Secret   string
}

type Daemon struct {
    Options *Options
    Front   *server.Server
    Backend *server.Server
    Storage *storage.Driver
}

type (
    CustomClaims struct {
        Device    string `json:"device"`
        Unique    string `json:"unique"`
        Timestamp int64  `json:"timestamp"`
        Addr      string `json:"addr"`
        jwt.StandardClaims
    }
)

func (d *Daemon) NewStorage() (err error) {
    conf := d.Options.Storage
    d.Storage, err = storage.New(conf.Driver, conf.Opts)
    return err
}

func (d *Daemon) NewFront() (err error) {
    d.Front = &server.Server{}
    Addr := fmt.Sprintf("%s:%s", d.Options.Host, d.Options.Port)
    if !d.Options.Security.TLS && !d.Options.Security.Verify {
        err = d.Front.Run(iris.Addr(Addr))
    } else {
        runner := iris.TLS(Addr, d.Options.Security.Cert, d.Options.Security.Key)
        err = d.Front.Run(runner)
    }
    if err == nil {
        d.Front.AddRouter(RouteToken.NewRouter(d))
    }
    return err
}

func (d *Daemon) NewBackend() (err error) {
    d.Backend = &server.Server{}
    l, err := netutil.UNIX("/tmpl/srv.sock", 0666)
    if err != nil {
        return err
    }
    err = d.Backend.Run(iris.Listener(l))
    if err == nil {
        // todo add backend router
    }
    return err
}

func (d *Daemon) Generate(req token.GenerateRequest) (string, error) {
    Claims := CustomClaims{
        req.Device,
        req.Unique,
        time.Now().Unix(),
        req.Addr,
        jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Second * TOKEN_TTL).Unix(),
            Issuer: "shuc324@gmail.com",
        },
    }
    Token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
    if Signed, err := Token.SignedString([]byte(d.Options.Secret)); err != nil {
        return "", err
    } else {
        err := (*d.Storage).LKeep(req.Unique, Signed, ALLOW_LOGIN_NUM, TOKEN_TTL)
        if err != nil {
            return "", err
        }
        return Signed, err
    }
}

func (d *Daemon) Auth(req token.AuthRequest) (interface{}, error) {
    Token, err := jwt.ParseWithClaims(
        req.JsonWebToken,
        &CustomClaims{},
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
            }
            return []byte(d.Options.Secret), nil
        })
    if err == nil && Token.Valid {
        if Claims, ok := Token.Claims.(*CustomClaims); ok {
            if (*d.Storage).LExist(Claims.Unique, req.JsonWebToken) {
                return Claims, nil
            }
        }
    }
    return nil, err
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
                return
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
    var err error;
    
    if args.Version == true {
        fmt.Printf("JwtAuth version %s.\n", VERSION)
        os.Exit(0)
    }
    
    Process := NewDaemon(args.Daemon, args)
    
    if Process.Options.Secret == "" {
        logrus.Error("please specify the key.")
        os.Exit(0)
    }
    
    if err = Process.NewStorage(); err != nil {
        logrus.Error(err)
        os.Exit(0)
    }
    
    if err = Process.NewFront(); err != nil {
        logrus.Error(err)
        os.Exit(0)
    }
    
    if err = Process.NewBackend(); err != nil {
        logrus.Error(err)
        os.Exit(0)
    }
}