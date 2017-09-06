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
    "github.com/BluePecker/JwtAuth/server/router"
    _ "github.com/BluePecker/JwtAuth/storage/redis"
    //_ "github.com/BluePecker/JwtAuth/storage/ram"
    "github.com/dgrijalva/jwt-go"
    RouteToken "github.com/BluePecker/JwtAuth/server/router/token"
)

const (
    TOKEN_TTL = 2 * 3600
    
    VERSION = "1.0.0"
    
    ALLOW_LOGIN_NUM = 3
)

type Storage struct {
    Driver string
    Opts    string
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
    Server  *server.Server
    Storage storage.Driver
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

func (d *Daemon) NewStorage() (*storage.Driver, error) {
    driver, err := storage.New(d.Options.Storage.Driver, d.Options.Storage.Opts)
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
    
    if d.Options.Security.Verify || d.Options.Security.TLS {
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
        err := d.Storage.LKeep(req.Unique, Signed, ALLOW_LOGIN_NUM, TOKEN_TTL)
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
            if d.Storage.LExist(Claims.Unique, req.JsonWebToken) {
                return Claims, nil
            }
        }
    }
    return nil, err
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
            logrus.Error(err)
            os.Exit(0)
        }
        d, err := dCtx.Reborn()
        if err != nil {
            if err == daemon.ErrWouldBlock {
                fmt.Println("daemon already exists.")
            } else {
                fmt.Println("Unable to run: ", err)
            }
            os.Exit(0)
        }
        if d != nil {
            return
        }
        defer dCtx.Release()
    }
    
    Daemon := &Daemon{
        Options: &args,
    }
    
    Storage, err := Daemon.NewStorage()
    if err != nil {
        logrus.Error(err)
        os.Exit(0)
    }
    Daemon.Storage = *Storage
    
    if Daemon.Options.Secret == "" {
        logrus.Error("please specify the key.")
        os.Exit(0)
    }
    
    Daemon.addRouter(RouteToken.NewRouter(Daemon))
    Daemon.Listen()
}