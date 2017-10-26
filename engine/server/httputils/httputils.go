package httputils

import (
	"github.com/kataras/iris/context"
	"github.com/kataras/iris"
	"github.com/BluePecker/JwtAuth/engine/server/parameter"
)

func Success(ctx context.Context, data interface{}) error {
	_, err := ctx.JSON(parameter.Response{
		Code:    iris.StatusOK,
		Data:    data,
		Message: "winner winner,chicken dinner.",
	})
	return err
}

func Failure(ctx context.Context, message string) error {
	_, err := ctx.JSON(parameter.Response{
		Code:    iris.StatusBadRequest,
		Data:    map[string]interface{}{},
		Message: message,
	})
	return err
}
