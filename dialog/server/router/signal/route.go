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

func (r *Route) Routes(api *iris.Application) {
	api.Get("/"+router.Version+"/stop", r.stop)
}

func NewRoute(b Backend) *Route {
	return &Route{backend: b}
}
