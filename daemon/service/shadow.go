package service

import (
	"context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/service/router"
	"github.com/BluePecker/JwtAuth/dialog/server"
)

type (
	Shadow struct {
		Routes    []router.Router
		WebServer *server.WebServer
	}
)

func (s *Shadow) New(ch chan struct{}, runner iris.Runner, configurator ... iris.Configurator) error {
	s.WebServer = &server.WebServer{Engine: iris.New()}

	for _, route := range s.Routes {
		s.WebServer.AddRouter(route)
	}
	configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
	configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
		DisableStartupLog: true,
	}))
	return s.WebServer.Run(runner, configurator...)
}

func (s *Shadow) Shutdown() error {
	return s.WebServer.Engine.Shutdown(context.TODO())
}
