package token

import (
	"github.com/kataras/iris/context"
	"github.com/BluePecker/JwtAuth/engine/server/parameter/jwt/request"
	"github.com/BluePecker/JwtAuth/engine/server/httputils"
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

func (r *Route) kick(ctx context.Context) {
	req := &request.Kick{}
	if err := ctx.ReadJSON(req); err != nil {
		httputils.Failure(ctx, err.Error())
	} else {
		if err := r.backend.Kick(*req); err != nil {
			httputils.Failure(ctx, err.Error())
		} else {
			httputils.Success(ctx, "ok")
		}
	}
}
