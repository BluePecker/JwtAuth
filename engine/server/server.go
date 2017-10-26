package server

import (
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/engine/server/router"
	"github.com/BluePecker/JwtAuth/engine/server/middleware"
)

type (
	WebServer struct {
		Engine *iris.Application
		routes []router.Route
	}
)

func (w *WebServer) AddRoute(routes ... router.Route) {
	for _, route := range routes {
		w.routes = append(w.routes, route)
	}
}

func (w *WebServer) Run(runner iris.Runner, configurator ... iris.Configurator) error {
	for _, ware := range middleware.Provider {
		w.Engine.Use(ware)
	}
	for _, route := range w.routes {
		route.Routes(w.Engine)
	}
	return w.Engine.Run(runner, configurator...)
}
