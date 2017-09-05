package token

import (
    "github.com/kataras/iris"
)

type Router struct {
    backend Backend
}

func (r *Router) Routes(server *iris.Application) {
    jwtRoutes := server.Party("/v1/token")
    {
        jwtRoutes.Post("/auth", r.auth)
        
        jwtRoutes.Post("/generate", r.generate)
    }
}

func NewRouter(backend Backend) *Router {
    return &Router{backend: backend}
}