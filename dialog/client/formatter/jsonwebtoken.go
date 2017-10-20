package formatter

import (
	"github.com/BluePecker/JwtAuth/dialog/client/formatter/context"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/response"
)

type (
	JsonWebToken struct {
		context.BaseSubjectContext
		truncate bool
		jwt      response.JsonWebToken
	}

	JsonWebTokenContext struct {
		context.Context
		JsonWebTokens []response.JsonWebToken
	}
)

func (c JsonWebTokenContext) Write() {
}
