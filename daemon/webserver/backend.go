package webserver

import (
	"context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/service/router"
	"github.com/BluePecker/JwtAuth/dialog/server"
)

type (
	Backend struct {
		Routes    []router.Router
		WebServer *server.WebServer
	}
)

func (b *Backend) New(ch chan struct{}, runner iris.Runner, configurator ... iris.Configurator) error {
	b.WebServer = &server.WebServer{Engine: iris.New()}

	for _, route := range b.Routes {
		b.WebServer.AddRouter(route)
	}
	configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
	configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}))
	return b.WebServer.Run(runner, configurator...)
}

func (b *Backend) Shutdown() error {
	return b.WebServer.Engine.Shutdown(context.TODO())
}
