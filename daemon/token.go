package daemon

import (
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/token/request"
	coderQ "github.com/BluePecker/JwtAuth/dialog/server/parameter/coder/request"
	"github.com/BluePecker/JwtAuth/dialog/server/parameter/token/response"
	"github.com/BluePecker/JwtAuth/daemon/coder"
)

func (d *Daemon) List(req request.List) ([]response.Token, error) {
	if sings, err := (*d.Cache).LRange(req.Unique, 0, coder.LoginNum); err != nil {
		return nil, err
	} else {
		ttl := (*d.Cache).TTL(req.Unique)
		tokens := []response.Token{}
		for _, singed := range sings {
			if token, err := coder.Decode(coderQ.Decode{
				JsonWebToken: singed,
			}, (*d.Options).Secret); err == nil && token != nil {
				if claims, ok := token.Claims.(*coder.CustomClaim); ok {
					tokens = append(tokens, response.Token{
						Singed: singed,
						TTL:    ttl,
						Addr:   claims.Addr,
						Device: claims.Device,
					})
				}
			}
		}
		return tokens, nil
	}
}
