package reply

import (
    "github.com/kataras/iris/context"
    "github.com/kataras/iris"
)

func Success(ctx context.Context, data interface{}) error {
    _, err := ctx.JSON(map[string]interface{}{
        "code": iris.StatusOK,
        "data": data,
        "message": "winner winner,chicken dinner.",
    })
    return err
}