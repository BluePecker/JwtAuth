package daemon

import (
    "github.com/BluePecker/JwtAuth/daemon/server"
    "github.com/BluePecker/JwtAuth/api/router"
    "github.com/BluePecker/JwtAuth/api/router/token"
    "fmt"
    "github.com/kataras/iris"
    "github.com/kataras/iris/core/netutil"
)

func (d *Daemon) Shadow(ch chan struct{}) error {
    d.shadow = &server.Shadow{}
    Listener, err := netutil.UNIX(d.Options.SockFile, 0666)
    if err != nil {
        return nil
    }
    return d.shadow.New(ch, iris.Listener(Listener))
}

func (d *Daemon) Rosiness(ch chan struct{}) error {
    d.rosiness = &server.Rosiness{
        Routes:[]router.Router{token.NewRouter(d)},
    }
    Addr := fmt.Sprintf("%s:%d", d.Options.Host, d.Options.Port)
    if d.Options.TLS.Cert != "" && d.Options.TLS.Key != "" {
        return d.rosiness.New(ch, iris.Addr(Addr))
    }
    runner := iris.TLS(Addr, d.Options.TLS.Cert, d.Options.TLS.Key)
    return d.rosiness.New(ch, runner)
}