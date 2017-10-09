package coder

import (
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
)

type (
	Route struct {
		backend Backend
	}
)

func (r *Route) Routes(server *iris.Application) {
	Route := server.Party("/" + router.Version + "/coder")
	{
		Route.Post("/decode", r.decode)

		Route.Post("/encode", r.encode)
	}
}

func NewRoute(backend Backend) *Route {
	return &Route{backend: backend}
}
