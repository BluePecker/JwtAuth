package server

import (
    "context"
    "github.com/kataras/iris"
    "github.com/BluePecker/JwtAuth/server/router"
    "github.com/BluePecker/JwtAuth/server"
)

type (
    Shadow struct {
        Routes  []router.Router
        Service *server.Server
    }
)

func (r *Shadow) New(runner iris.Runner, configurator iris.Configurator) error {
    r.Service = &server.Server{App: iris.New()}
    
    for _, route := range r.Routes {
        r.Service.AddRouter(route)
    }
    return r.Service.Run(runner, iris.WithConfiguration(iris.Configuration{
        DisableStartupLog: true,
    }), iris.WithoutServerError(iris.ErrServerClosed), configurator)
}

func (r *Shadow) Shutdown() error {
    return r.Service.App.Shutdown(context.TODO())
}