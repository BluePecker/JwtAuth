package middleware

import (
	"time"
	"github.com/kataras/iris/context"
	"github.com/Sirupsen/logrus"
	"strconv"
)

func init() {
	Register(func(ctx context.Context) {
		start := time.Now()
		logrus.Info("test")
		ctx.Next()
		logrus.Infof("%v %4v %s %s %s", strconv.Itoa(ctx.GetStatusCode()), time.Now().Sub(start), ctx.RemoteAddr(), ctx.Method(), ctx.Path())
	})
}
