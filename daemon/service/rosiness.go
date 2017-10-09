package service

import (
	"context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/service/router"
	"github.com/BluePecker/JwtAuth/dialog/server"
)

type (
	Rosiness struct {
		Routes    []router.Router
		WebServer *server.WebServer
	}
)

func (r *Rosiness) New(ch chan struct{}, runner iris.Runner, configurator ... iris.Configurator) error {
	r.WebServer = &server.WebServer{Engine: iris.New()}
	for _, route := range r.Routes {
		r.WebServer.AddRouter(route)
	}
	go func() {
		if _, ok := <-ch; ok {
			r.Shutdown()
		}
	}()
	configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
	configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}))
	return r.WebServer.Run(runner, configurator...)
}

func (r *Rosiness) Shutdown() error {
	return r.WebServer.Engine.Shutdown(context.TODO())
}
