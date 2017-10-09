package daemon

import (
	"github.com/BluePecker/JwtAuth/service/router"
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/netutil"
	"github.com/BluePecker/JwtAuth/service/router/signal"
	"github.com/BluePecker/JwtAuth/daemon/webserver"
)

func (d *Daemon) Backend(ch chan struct{}) error {
	d.backend = &webserver.Backend{
		Routes: []router.Router{signal.NewRouter(d)},
	}
	Listener, err := netutil.UNIX(d.Options.SockFile, 0666)
	if err != nil {
		return nil
	}
	return d.backend.New(ch, iris.Listener(Listener))
}

func (d *Daemon) Front(ch chan struct{}) error {
	d.front = &webserver.Front{
		Routes: []router.Router{},
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
