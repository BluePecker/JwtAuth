package version

import (
    "github.com/kataras/iris"
)

type Router struct {
    backend Backend
}

func (r *Router) Routes(server *iris.Application) {
    server.Get("/v1/version", r.version)
}

func NewRouter(backend Backend) *Router {
    return &Router{backend: backend}
}
