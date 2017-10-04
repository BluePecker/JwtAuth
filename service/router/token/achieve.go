package token

import (
    "github.com/kataras/iris/context"
    "github.com/BluePecker/JwtAuth/service/types/token"
    Response "github.com/BluePecker/JwtAuth/service/reply"
)

func (r *Router) auth(ctx context.Context) {
    req := &token.AuthRequest{}
    if err := ctx.ReadJSON(req); err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    claims, err := r.backend.Auth(*req)
    if err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    Response.Success(ctx, claims)
}

func (r *Router) generate(ctx context.Context) {
    req := &token.GenerateRequest{}
    if err := ctx.ReadJSON(req); err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    jwt, err := r.backend.Generate(*req)
    if err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    Response.Success(ctx, map[string]interface{}{
        "token": jwt,
    })
}