package token

import (
	"github.com/BluePecker/JwtAuth/engine/server/parameter/jwt/request"
	"github.com/BluePecker/JwtAuth/engine/server/parameter/jwt/response"
)

type (
	Backend interface {
		Kick(req request.Kick) error

		List(req request.List) ([]response.JsonWebToken, error)
	}
)
