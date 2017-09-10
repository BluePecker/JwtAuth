package server

import (
    "context"
    "github.com/kataras/iris"
    "github.com/BluePecker/JwtAuth/server/router"
    "github.com/BluePecker/JwtAuth/server"
)

type (
    Rosiness struct {
        Routes  []router.Router
        Service *server.Server
    }
)

func (s *Rosiness) New(runner iris.Runner, configurator iris.Configurator) error {
    s.Service = &server.Server{App: iris.New()}
    for _, route := range s.Routes {
        s.Service.AddRouter(route)
    }
    return s.Service.Run(runner, configurator)
}

func (s *Rosiness) Shutdown() error {
    return s.Service.App.Shutdown(context.TODO())
}