package daemon

import (
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/request"
	coderQ "github.com/BluePecker/JwtAuth/dialog/server/parameter/coder/request"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/jwt/response"
	"github.com/BluePecker/JwtAuth/daemon/coder"
	"github.com/Sirupsen/logrus"
)

func (d *Daemon) List(req request.List) ([]response.JsonWebToken, error) {
	if sings, err := (*d.Cache).LRange(req.Unique, 0, coder.LoginNum); err != nil {
		return nil, err
	} else {
		ttl := (*d.Cache).TTL(req.Unique)
		tokens := []response.JsonWebToken{}
		for _, singed := range sings {
			if token, err := coder.Decode(coderQ.Decode{
				JsonWebToken: singed,
			}, (*d.Options).Secret); err == nil {
				if claims, ok := token.Claims.(*coder.CustomClaim); ok {
					tokens = append(tokens, response.JsonWebToken{
						Singed: singed,
						TTL:    ttl,
						Addr:   claims.Addr,
						Device: claims.Device,
					})
				}
			} else {
				logrus.Error(err)
			}
		}
		return tokens, nil
	}
}
