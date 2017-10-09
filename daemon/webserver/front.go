package webserver

import (
	"context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/dialog/server"
)

type (
	Front Web
)

func (f *Front) New(ch chan struct{}, runner iris.Runner, configurator ... iris.Configurator) error {
	f.App = &server.WebServer{Engine: iris.New()}
	for _, route := range f.Routes {
		f.App.AddRouter(route)
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
	return f.App.Run(runner, configurator...)
}

func (f *Front) Shutdown() error {
	return f.App.Engine.Shutdown(context.TODO())
}
