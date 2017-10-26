package daemon

import (
	"github.com/BluePecker/JwtAuth/engine/server/parameter/coder/request"
	"github.com/BluePecker/JwtAuth/daemon/coder"
)

func (d *Daemon) Decode(req request.Decode) (*coder.CustomClaim, error) {
	token, err := coder.Decode(req, d.Options.Secret)
	if err == nil {
		if claims, ok := token.Claims.(*coder.CustomClaim); ok {
			cache, ttl, err := (*d.Cache).HGet(claims.Unique, claims.Device)
			if ttl >= 0 && cache == req.JsonWebToken {
				return claims, nil
			} else {
				return nil, err
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
	err = (*d.Cache).HSet(req.Unique, req.Device, singed, coder.LoginNum, coder.TokenTTL)
	if err != nil {
		return "", err
	}
	return singed, nil
}
