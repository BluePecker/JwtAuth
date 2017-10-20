package token

import (
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/request"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/response"
)

type (
	Backend interface {
		List(req request.List) ([]response.JsonWebToken, error)
	}
)
