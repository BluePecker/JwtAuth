package daemon

import (
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"github.com/BluePecker/JwtAuth/service/types/token"
	"time"
)

const (
	AllowLoginNum = 3
	TokenTTL      = 2 * 3600
)

type (
	CustomClaim struct {
		Device    string `json:"device"`
		Unique    string `json:"unique"`
		Timestamp int64  `json:"timestamp"`
		Addr      string `json:"addr"`
		jwt.StandardClaims
	}
)

func (d *Daemon) Generate(req token.GenerateRequest) (string, error) {
	Claims := CustomClaim{
		req.Device,
		req.Unique,
		time.Now().Unix(),
		req.Addr,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * TokenTTL).Unix(),
			Issuer:    "shuc324@gmail.com",
		},
	}
	Token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)
	if Signed, err := Token.SignedString([]byte(d.Options.Secret)); err != nil {
		return "", err
	} else {
		err := (*d.StorageE).LKeep(req.Unique, Signed, AllowLoginNum, TokenTTL)
		if err != nil {
			return "", err
		}
		return Signed, err
	}
}

func (d *Daemon) Auth(req token.AuthRequest) (interface{}, error) {
	Token, err := jwt.ParseWithClaims(
		req.JsonWebToken,
		&CustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%v", token.Header["alg"])
			}
			return []byte(d.Options.Secret), nil
		})
	if err == nil && Token.Valid {
		if Claims, ok := Token.Claims.(*CustomClaim); ok {
			if (*d.StorageE).LExist(Claims.Unique, req.JsonWebToken) {
				return Claims, nil
			}
		}
	}
	return nil, err
}
