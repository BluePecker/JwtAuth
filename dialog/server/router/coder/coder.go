package coder

import (
	"github.com/kataras/iris/context"
	request "github.com/BluePecker/JwtAuth/dialog/server/parameter/coder"
	"github.com/BluePecker/JwtAuth/dialog/server/httputils"
)

func (r *Route) decode(ctx context.Context) {
	req := &request.Decode{}
	if err := ctx.ReadJSON(req); err != nil {
		httputils.Failure(ctx, err.Error())
		return
	}
	claims, err := r.backend.Decode(*req)
	if err != nil {
		httputils.Failure(ctx, err.Error())
		return
	}
	httputils.Success(ctx, claims)
}

func (r *Route) encode(ctx context.Context) {
	req := &request.Encode{}
	if err := ctx.ReadJSON(req); err != nil {
		httputils.Failure(ctx, err.Error())
		return
	}
	jwt, err := r.backend.Encode(*req)
	if err != nil {
		httputils.Failure(ctx, err.Error())
		return
	}
	httputils.Success(ctx, map[string]interface{}{
		"token": jwt,
	})
}
