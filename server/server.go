package server

import (
    "strconv"
    "time"
    "github.com/kataras/iris"
    "github.com/BluePecker/JwtAuth/server/router"
    "github.com/Sirupsen/logrus"
    "github.com/kataras/iris/context"
)

type (
    Server struct {
        App *iris.Application
    }
)

func (Api *Server) AddRouter(routers... router.Router) {
    for _, route := range routers {
        route.Routes(Api.App)
    }
}

func (Api *Server) Run(runner iris.Runner) error {
    Api.App = iris.New()
    Options := iris.WithConfiguration(iris.Configuration{
        //DisableStartupLog: true,
    })
    Api.App.Use(func(ctx context.Context) {
        start := time.Now()
        ctx.Next()
        logrus.Infof("%v %4v %s %s %s", strconv.Itoa(ctx.GetStatusCode()), time.Now().Sub(start), ctx.RemoteAddr(), ctx.Method(), ctx.Path())
    })
    return Api.App.Run(runner, Options)
}