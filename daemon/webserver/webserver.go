package webserver

import (
	"context"
	"github.com/BluePecker/JwtAuth/dialog/server"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
	"github.com/kataras/iris"
)

type (
	Web struct {
		App *server.WebServer
	}
)

func (w *Web) ShutDown() error {
	return w.App.Engine.Shutdown(context.TODO())
}

func (w *Web) Listen(runner iris.Runner, ch chan struct{}, route ... router.Route) error {
	w.App = &server.WebServer{Engine: iris.New()}

	w.App.AddRoute(route...)

	if ch != nil {
		go func() {
			if _, ok := <-ch; ok {
				w.ShutDown()
			}
		}()
	}

	return w.App.Run(runner, iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}), iris.WithoutServerError(iris.ErrServerClosed))
}
