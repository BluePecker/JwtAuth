package webserver

import (
	"context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/dialog/server"
)

type (
	Front Web
)

func (f *Front) New(ch chan struct{}, runner iris.Runner) error {
	f.App = &server.WebServer{Engine: iris.New()}
	f.App.AddRoute(f.Routes...)
	go func() {
		if _, ok := <-ch; ok {
			f.Shutdown()
		}
	}()
	var configurator []iris.Configurator
	configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
	configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}))
	return f.App.Run(runner, configurator...)
}

func (f *Front) Shutdown() error {
	return f.App.Engine.Shutdown(context.TODO())
}
