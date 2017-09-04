package token

import (
    "github.com/kataras/iris/context"
    "github.com/kataras/iris"
)

type (
    User struct {
        _id string
    }
)

func (auth *authRouter) generate(ctx context.Context) {
    user := &User{}
    if err := ctx.ReadJSON(user); err != nil {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusBadRequest,
            "data": map[string]interface{}{},
            "message": err.Error(),
        })
        return
    }
    Token, err := auth.standard.Generate(user._id)
    if err != nil {
        ctx.JSON(map[string]interface{}{
            "code": iris.StatusBadRequest,
            "data": map[string]interface{}{},
            "message": err.Error(),
        })
        return
    }
    
    ctx.JSON(map[string]interface{}{
        "code": iris.StatusOK,
        "data": Token,
        "message": "winner winner,chicken dinner.",
    })
    return
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