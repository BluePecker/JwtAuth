package token

import (
    "github.com/kataras/iris"
)

type Router struct {
    backend Backend
}

func (r *Router) Routes(server *iris.Application) {
    Route := server.Party("/v1/token")
    {
        Route.Post("/auth", r.auth)
        
        Route.Post("/generate", r.generate)
    }
}

func NewRouter(backend Backend) *Router {
    return &Router{backend: backend}
}