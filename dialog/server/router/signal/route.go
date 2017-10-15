package signal

import (
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
)

type (
	Route struct {
		backend Backend
	}
)

func (r *Route) Routes(app *iris.Application) {
	Route := app.Party("/" + router.Version + "/signal")
	{
		Route.Get("/stop", r.stop)
	}
}

func NewRoute(b Backend) *Route {
	return &Route{backend: b}
}
