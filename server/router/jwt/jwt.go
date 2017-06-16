package jwt

import (
    "github.com/kataras/iris"
)

type jwtRouter struct {
    standard Standard
}

func (r *jwtRouter) Routes(server *iris.Application) {
    
    jwtRoutes := server.Party("/jwt")
    {
        jwtRoutes.Post("/generate", r.generate)
        
        jwtRoutes.Post("/auth", r.auth)
        
        jwtRoutes.Put("/upgrade", r.upgrade)
    }
}

func NewRouter(standard Standard) *jwtRouter {
    return &jwtRouter{
        standard: standard,
    }
}