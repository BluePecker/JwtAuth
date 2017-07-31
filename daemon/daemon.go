package daemon

import (
    "github.com/sevlyar/go-daemon"
    "github.com/Sirupsen/logrus"
    "github.com/BluePecker/JwtAuth/server"
    "github.com/BluePecker/JwtAuth/server/router/jwt"
    "github.com/BluePecker/JwtAuth/storage"
    "reflect"
    "github.com/kataras/iris/core/router"
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

type Option struct {
    PidFile  string
    LogFile  string
    Port     int
    Host     string
    Daemon   bool
    Security Security
    Storage  Storage
}

type Daemon struct {
    opt     *Option
    server  *server.Server
    storage *storage.Driver
}

func (d *Daemon) storageConf(p2 *storage.Option) {
    p1 := d.opt.Storage
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

func (d *Daemon) initStorage() {
    option := &storage.Option{}
    driver := d.opt.Storage.Driver
    d.storageConf(option)
    
    driver, err := storage.New(driver, option)
    if err != nil {
        logrus.Error(err)
        return
    }
    d.storage = driver
}

func (d *Daemon) initServer() {
    d.server = &server.Server{}
}

func (d *Daemon) listen() {
    if d.server == nil {
        d.initServer()
    }
    
    d.server.Accept(server.Options{
        Host: d.opt.Host,
        Port: d.opt.Port,
        Tls: &server.TLS{
            Key: d.opt.Security.Key,
            Cert: d.opt.Security.Cert,
        },
    })
}

func (d *Daemon) addRouter(routers... router.Router) {
    if d.server == nil {
        d.initServer()
    }
    d.server.AddRouter(routers)
}

func NewStart(opt Option) {
    if (opt.Daemon == true) {
        dCtx := daemon.Context{
            PidFileName: opt.PidFile,
            PidFilePerm: 0644,
            LogFilePerm: 0640,
            Umask:       027,
            WorkDir:     "/",
            LogFileName: opt.LogFile,
        }
        
        defer dCtx.Release()
        
        if child, err := dCtx.Reborn(); err != nil {
            logrus.Fatal(err)
        } else if child != nil {
            return
        }
    }
    
    jwtPro := &Daemon{
        opt: &opt,
    }
    
    jwtPro.initStorage()
    
    jwtPro.addRouter(jwt.NewRouter(nil))
    jwtPro.listen()
}