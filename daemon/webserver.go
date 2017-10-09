package daemon

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/netutil"
	"github.com/BluePecker/JwtAuth/daemon/webserver"
	"github.com/BluePecker/JwtAuth/dialog/server/router/signal"
	"github.com/BluePecker/JwtAuth/dialog/server/router/coder"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
)

func (d *Daemon) Backend(ch chan struct{}) error {
	d.backend = &webserver.Backend{
		Routes: []router.Route{signal.NewRoute(d), coder.NewRoute(d)},
	}
	Listener, err := netutil.UNIX(d.Options.SockFile, 0666)
	if err != nil {
		return nil
	}
	return d.backend.New(ch, iris.Listener(Listener))
}

func (d *Daemon) Front(ch chan struct{}) error {
	d.front = &webserver.Front{
		Routes: []router.Route{},
	}
	Addr := fmt.Sprintf("%s:%d", d.Options.Host, d.Options.Port)
	var runner iris.Runner
	if d.Options.TLS.Cert != "" && d.Options.TLS.Key != "" {
		runner = iris.TLS(Addr, d.Options.TLS.Cert, d.Options.TLS.Key)
	} else {
		runner = iris.Addr(Addr)
	}
	return d.front.New(ch, runner)
}
