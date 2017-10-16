package token

import (
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/token/request"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/token/response"
)

type (
	Backend interface {
		List(req request.List) ([]response.Token, error)
	}
)
