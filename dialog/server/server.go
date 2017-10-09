package server

import (
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
	"github.com/BluePecker/JwtAuth/dialog/server/middleware"
)

type (
	WebServer struct {
		Engine *iris.Application
	}
)

func (w *WebServer) AddRouter(routes ... router.Route) {
	for _, route := range routes {
		route.Routes(w.Engine)
	}
}

func (w *WebServer) Run(runner iris.Runner, configurator ... iris.Configurator) error {
	for _, ware := range middleware.Provider {
		w.Engine.Use(ware)
	}
	return w.Engine.Run(runner, configurator...)
}
