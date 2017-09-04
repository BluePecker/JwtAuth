package token

import (
    "github.com/kataras/iris"
)

type authRouter struct {
    standard Standard
}

func (r *authRouter) Routes(server *iris.Application) {
    jwtRoutes := server.Party("/v1/token")
    {
        jwtRoutes.Post("/generate", r.generate)
        
        jwtRoutes.Post("/auth", r.auth)
        
        jwtRoutes.Put("/upgrade", r.upgrade)
    }
}

func NewRouter(standard *Standard) *authRouter {
    return &authRouter{
        standard: standard,
    }
}