package reply

import (
    "github.com/kataras/iris/context"
    "github.com/kataras/iris"
)

func Failure(ctx context.Context, message string) error {
    _, err := ctx.JSON(map[string]interface{}{
        "code": iris.StatusBadRequest,
        "data": map[string]interface{}{},
        "message": message,
    })
    return err
}