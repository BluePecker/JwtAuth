package router

import "github.com/kataras/iris"

type Router interface {
    Routes(server *iris.Application)
}

