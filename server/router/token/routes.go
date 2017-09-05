package token

import (
    "github.com/kataras/iris/context"
    "github.com/BluePecker/JwtAuth/server/types/token"
    Response "github.com/BluePecker/JwtAuth/server/reply"
)

func (auth *Router) auth(ctx context.Context) {
    req := &token.AuthRequest{}
    if err := ctx.ReadJSON(req); err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    claims, err := auth.backend.Auth(*req)
    if err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    Response.Success(ctx, claims)
}

func (auth *Router) generate(ctx context.Context) {
    req := &token.GenerateRequest{}
    if err := ctx.ReadJSON(req); err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    jwt, err := auth.backend.Generate(*req)
    if err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    Response.Success(ctx, map[string]interface{}{
        "token": jwt,
    })
}