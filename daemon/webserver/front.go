package webserver

import (
	"context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/dialog/server"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
)

type (
	Front struct {
		Routes    []router.Route
		WebServer *server.WebServer
	}
)

func (f *Front) New(ch chan struct{}, runner iris.Runner, configurator ... iris.Configurator) error {
	f.WebServer = &server.WebServer{Engine: iris.New()}
	for _, route := range f.Routes {
		f.WebServer.AddRouter(route)
	}
	go func() {
		if _, ok := <-ch; ok {
			f.Shutdown()
		}
	}()
	configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
	configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}))
	return f.WebServer.Run(runner, configurator...)
}

func (f *Front) Shutdown() error {
	return f.WebServer.Engine.Shutdown(context.TODO())
}
