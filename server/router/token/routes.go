package token

import (
    "github.com/kataras/iris/context"
    "github.com/kataras/iris"
)

type (
    Payload struct {
        Device  string `json:"device"`
        UserId  string `json:"user_id"`
        Address string `json:"address"`
    }
)

func (auth *authRouter) generate(ctx context.Context) {
    Payload := &Payload{}
    if err := ctx.ReadJSON(Payload); err != nil {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusBadRequest,
            "data": map[string]interface{}{},
            "message": err.Error(),
        })
        return
    }
    Token, err := auth.standard.Generate(Payload.UserId, Payload.Device, Payload.Address)
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
            "data": Token,
            "message": "winner winner,chicken dinner.",
        })
        return
    }
}

func (auth *authRouter) auth(ctx context.Context) {
    ctx.JSON(map[string]interface{}{
        "user_id": 10000,
    })
}

func (auth *authRouter) upgrade(ctx context.Context) {
    ctx.JSON(map[string]interface{}{
        "user_id": 10000,
    })
}