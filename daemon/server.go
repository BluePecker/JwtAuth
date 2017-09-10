package daemon

import (
    "github.com/BluePecker/JwtAuth/daemon/server"
    "github.com/BluePecker/JwtAuth/server/router"
    "github.com/BluePecker/JwtAuth/server/router/token"
    "fmt"
    "github.com/kataras/iris"
    "github.com/kataras/iris/core/netutil"
)

func (d *Daemon) Shadow() error {
    d.shadow = &server.Shadow{}
    Listener, err := netutil.UNIX(d.Options.SockFile, 0666)
    if err != nil {
        return nil
    }
    return d.shadow.New(iris.Listener(Listener))
}

func (d *Daemon) Rosiness() error {
    d.rosiness = &server.Rosiness{
        Routes:[]router.Router{token.NewRouter(d)},
    }
    Addr := fmt.Sprintf("%s:%d", d.Options.Host, d.Options.Port)
    if d.Options.TLS.Cert != "" && d.Options.TLS.Key != "" {
        return d.rosiness.New(iris.Addr(Addr))
    }
    runner := iris.TLS(Addr, d.Options.TLS.Cert, d.Options.TLS.Key)
    return d.rosiness.New(runner)
}