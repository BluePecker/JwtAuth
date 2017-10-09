package daemon

import (
	request "github.com/BluePecker/JwtAuth/dialog/server/parameter/coder"
	"github.com/BluePecker/JwtAuth/daemon/coder"
	"github.com/kataras/iris/core/errors"
)

func (d *Daemon) Decode(req request.Decode) (*coder.CustomClaim, error) {
	token, err := coder.Decode(req, d.Options.Secret)
	if err == nil {
		if claims, ok := token.Claims.(*coder.CustomClaim); ok {
			if !(*d.Cache).LExist(claims.Unique, req.JsonWebToken) {
				return nil, errors.New("illegal json-web-token")
			} else {
				return claims, nil
			}
		}
	}
	return nil, err
}

func (d *Daemon) Encode(req request.Encode) (string, error) {
	singed, err := coder.Encode(req, d.Options.Secret)
	if err != nil {
		return "", err
	}
	err = (*d.Cache).LKeep(req.Unique, singed, coder.LoginNum, coder.TokenTTL)
	if err != nil {
		return "", err
	}
	return singed, nil
}
