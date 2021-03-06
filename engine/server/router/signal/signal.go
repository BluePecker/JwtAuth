package signal

import (
	"github.com/kataras/iris/context"
	"github.com/BluePecker/JwtAuth/engine/server/httputils"
)

func (r *Route) stop(ctx context.Context) {
	httputils.Success(ctx, "ok")
	r.backend.Stop()
}
