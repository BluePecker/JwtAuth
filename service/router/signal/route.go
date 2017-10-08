package signal

import (
	"github.com/kataras/iris"
)

type Router struct {
	backend Backend
}

func (r *Router) Routes(server *iris.Application) {
	server.Get("/v1/stop", r.stop)
}

func NewRouter(backend Backend) *Router {
	return &Router{backend: backend}
}