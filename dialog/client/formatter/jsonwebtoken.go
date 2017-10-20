package formatter

import (
	"github.com/BluePecker/JwtAuth/dialog/client/formatter/context"
	"github.com/dgrijalva/jwt-go"
)

type (
	JsonWebTokenContext struct {
		context.Context
		JsonWebTokens []jwt.Token
	}
)

func (c JsonWebTokenContext) Write() {
}
