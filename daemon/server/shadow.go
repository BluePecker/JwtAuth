package server

import (
    "context"
    "github.com/kataras/iris"
    "github.com/BluePecker/JwtAuth/api"
    "github.com/BluePecker/JwtAuth/api/router"
)

type (
    Shadow struct {
        Routes  []router.Router
        Service *server.Server
    }
)

func (r *Shadow) New(ch chan struct{}, runner iris.Runner, configurator... iris.Configurator) error {
    r.Service = &server.Server{App: iris.New()}
    
    for _, route := range r.Routes {
        r.Service.AddRouter(route)
    }
    //go func() {
    //    if _, ok := <-ch; ok {
    //        r.Shutdown()
    //    }
    //}()
    configurator = append(configurator, iris.WithoutServerError(iris.ErrServerClosed))
    configurator = append(configurator, iris.WithConfiguration(iris.Configuration{
        DisableStartupLog: true,
    }))
    return r.Service.Run(runner, configurator...)
}

func (r *Shadow) Shutdown() error {
    return r.Service.App.Shutdown(context.TODO())
}