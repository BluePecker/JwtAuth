package version

import (
    "github.com/kataras/iris/context"
    Response "github.com/BluePecker/JwtAuth/service/reply"
)

func (r *Router) version(ctx context.Context) {
    version, err := r.backend.Version()
    if err != nil {
        Response.Failure(ctx, err.Error())
        return
    }
    Response.Success(ctx, version)
}