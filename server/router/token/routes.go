package token

import (
    "github.com/kataras/iris/context"
    "github.com/kataras/iris"
)

type (
    Generate struct {
        Device  string `json:"device"`
        UserId  string `json:"user_id"`
        Address string `json:"address"`
    }
    
    Auth struct {
        Token string `json:"token"`
    }
)

func (auth *authRouter) generate(ctx context.Context) {
    G := &Generate{}
    if err := ctx.ReadJSON(G); err != nil {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusBadRequest,
            "data": map[string]interface{}{},
            "message": err.Error(),
        })
        return
    }
    Token, err := auth.standard.Generate(G.UserId, G.Device, G.Address)
    if err != nil {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusBadRequest,
            "data": map[string]interface{}{},
            "message": err.Error(),
        })
        return
    } else {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusOK,
            "data": map[string]interface{}{
                "token": Token,
            },
            "message": "winner winner,chicken dinner.",
        })
        return
    }
}

func (auth *authRouter) auth(ctx context.Context) {
    A := &Auth{}
    if err := ctx.ReadJSON(A); err != nil {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusBadRequest,
            "data": map[string]interface{}{},
            "message": err.Error(),
        })
        return
    }
    UserId, err := auth.standard.Auth(A.Token)
    if err != nil {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusBadRequest,
            "data": map[string]interface{}{},
            "message": err.Error(),
        })
        return
    } else {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusOK,
            "data": map[string]interface{}{
                "user_id": UserId,
            },
            "message": "winner winner,chicken dinner.",
        })
        return
    }
}

func (auth *authRouter) upgrade(ctx context.Context) {
    ctx.JSON(map[string]interface{}{
        "user_id": 10000,
    })
}