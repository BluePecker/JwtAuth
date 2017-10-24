package daemon

import (
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/request"
	coderQ "github.com/BluePecker/JwtAuth/dialog/server/parameter/coder/request"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/response"
	"github.com/BluePecker/JwtAuth/daemon/coder"
	"github.com/Sirupsen/logrus"
)

func (d *Daemon) List(req request.List) ([]response.JsonWebToken, error) {
	var tokens []response.JsonWebToken
	if err := (*d.Cache).HScan(req.Unique, func(singed string, ttl float64) {
		if token, err := coder.Decode(coderQ.Decode{
			JsonWebToken: singed,
		}, (*d.Options).Secret); err != nil {
			logrus.Error(err)
		} else {
			if claims, ok := token.Claims.(*coder.CustomClaim); ok {
				tokens = append(tokens, response.JsonWebToken{
					Singed: singed,
					TTL:    ttl,
					Addr:   claims.Addr,
					Device: claims.Device,
				})
			}
		}
	}); err != nil {
		return tokens, err
	} else {
		return tokens, nil
	}
}
