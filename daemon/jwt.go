package daemon

import (
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/request"
	coderQ "github.com/BluePecker/JwtAuth/dialog/server/parameter/coder/request"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/response"
	"github.com/BluePecker/JwtAuth/daemon/coder"
	"github.com/Sirupsen/logrus"
)

func (d *Daemon) List(req request.List) ([]response.JsonWebToken, error) {
	if keys, err := (*d.Cache).HKeys(req.Unique); err != nil {
		return nil, err
	} else {
		logrus.Info(keys)
		var tokens []response.JsonWebToken
		for _, k := range keys {
			if singed, ttl, err := (*d.Cache).HGetString(req.Unique, k); err == nil {
				logrus.Info(ttl, singed)
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
			}
		}
		return tokens, nil
	}
}
