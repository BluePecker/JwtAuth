package webserver

import (
	"github.com/BluePecker/JwtAuth/dialog/server"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
)

type (
	Web struct {
		Routes []router.Route
		App    *server.WebServer
	}
)

func (w *Web) AddRoute(route ... router.Route) {
	w.Routes = append(w.Routes, route...)
}
