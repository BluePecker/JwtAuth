package token

import (
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/engine/server/router"
)

type Route struct {
	backend Backend
}

func (r *Route) Routes(app *iris.Application) {
	Route := app.Party("/" + router.Version + "/token")
	{
		Route.Post("/list", r.list)
	}
}

func NewRoute(backend Backend) *Route {
	return &Route{backend: backend}
}
