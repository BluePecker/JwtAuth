package daemon

import (
    "github.com/BluePecker/JwtAuth/daemon/service"
    "github.com/BluePecker/JwtAuth/service/router"
    "github.com/BluePecker/JwtAuth/service/router/token"
    "fmt"
    "github.com/kataras/iris"
    "github.com/kataras/iris/core/netutil"
    "github.com/BluePecker/JwtAuth/service/router/version"
)

func (d *Daemon) Shadow(ch chan struct{}) error {
    d.shadow = &service.Shadow{
        Routes: []router.Router(version.NewRouter(d)),
    }
    Listener, err := netutil.UNIX(d.Options.SockFile, 0666)
    if err != nil {
        return nil
    }
    return d.shadow.New(ch, iris.Listener(Listener))
}

func (d *Daemon) Rosiness(ch chan struct{}) error {
    d.rosiness = &service.Rosiness{
        Routes:[]router.Router{token.NewRouter(d)},
    }
    Addr := fmt.Sprintf("%s:%d", d.Options.Host, d.Options.Port)
    var runner iris.Runner
    if d.Options.TLS.Cert != "" && d.Options.TLS.Key != "" {
        runner = iris.TLS(Addr, d.Options.TLS.Cert, d.Options.TLS.Key)
    } else {
        runner = iris.Addr(Addr)
    }
    return d.rosiness.New(ch, runner)
}