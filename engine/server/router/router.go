package router

import "github.com/kataras/iris"

const (
	Version = "v1.0"
)

type (
	Route interface {
		Routes(engine *iris.Application)
	}
)
