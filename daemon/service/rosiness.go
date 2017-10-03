package service

import (
    "context"
    "github.com/kataras/iris"
    "github.com/BluePecker/JwtAuth/service/router"
    "github.com/BluePecker/JwtAuth/service"
)

type (
    Rosiness struct {
        Routes  []router.Router
        Service *service.Server
    }
)

func (s *Rosiness) New(ch chan struct{}, runner iris.Runner, configurator... iris.Configurator) error {
    s.Service = &service.Server{App: iris.New()}
    for _, route := range s.Routes {
        s.Service.AddRouter(route)
    }
    go func() {
        if _, ok := <-ch; ok {
            s.Shutdown()
        }
    }()
    configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
    configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
        DisableStartupLog: true,
    }))
    return s.Service.Run(runner, configurator...)
}

func (s *Rosiness) Shutdown() error {
    return s.Service.App.Shutdown(context.TODO())
}