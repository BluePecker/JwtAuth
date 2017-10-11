package webserver

import (
	"context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/dialog/server"
)

type (
	Backend Web
)

func (b *Backend) New(ch chan struct{}, runner iris.Runner) error {
	b.App = &server.WebServer{Engine: iris.New()}
	b.App.AddRoute(b.Routes...)

	var configurator []iris.Configurator
	configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
	configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}))
	return b.App.Run(runner, configurator...)
}

func (b *Backend) Shutdown() error {
	return b.App.Engine.Shutdown(context.TODO())
}
