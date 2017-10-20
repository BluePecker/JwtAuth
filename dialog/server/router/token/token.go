package token

import (
	"github.com/kataras/iris/context"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/request"
	"github.com/BluePecker/JwtAuth/dialog/server/httputils"
)

func (r *Route) list(ctx context.Context) {
	req := &request.List{}
	if err := ctx.ReadJSON(req); err != nil {
		httputils.Failure(ctx, err.Error())
		return
	}
	list, err := r.backend.List(*req)
	if err != nil {
		httputils.Failure(ctx, err.Error())
		return
	}
	httputils.Success(ctx, list)
}
