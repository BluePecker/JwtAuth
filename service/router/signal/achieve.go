package signal

import (
	"github.com/kataras/iris/context"
	Response "github.com/BluePecker/JwtAuth/service/reply"
)

func (r *Router) stop(ctx context.Context) {
	Response.Success(ctx, "ok")
	r.backend.Stop()
}
