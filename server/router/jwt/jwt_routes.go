package jwt

import "github.com/kataras/iris/context"

func (a *jwtRouter) generate(ctx context.Context) {
    ctx.JSON(map[string]interface{}{
        "user_id": 10000,
    })
}

func (a *jwtRouter) auth(ctx context.Context) {
    ctx.JSON(map[string]interface{}{
        "user_id": 10000,
    })
}

func (a *jwtRouter) upgrade(ctx context.Context) {
    ctx.JSON(map[string]interface{}{
        "user_id": 10000,
    })
}