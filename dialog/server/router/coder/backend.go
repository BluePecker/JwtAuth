package coder

import (
	request "github.com/BluePecker/JwtAuth/dialog/server/parameter/coder"
	"github.com/BluePecker/JwtAuth/daemon/coder"
)

type Backend interface {
	Decode(req request.Decode) (*coder.CustomClaim, error)

	Encode(req request.Encode) (string, error)
}
