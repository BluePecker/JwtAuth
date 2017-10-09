package webserver

import (
	"github.com/BluePecker/JwtAuth/dialog/server"
	"github.com/BluePecker/JwtAuth/dialog/server/router"
)

type (
	Web struct {
		Routes []router.Route
		App    *server.WebServer
	}
)
