package coder

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"fmt"
	"github.com/BluePecker/JwtAuth/engine/server/parameter/coder/request"
)

const (
	LoginNum = 3
	TokenTTL = 2 * 3600
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

func Decode(req request.Decode, secret string) (*jwt.Token, error) {
	Token, err := jwt.ParseWithClaims(
		req.JsonWebToken,
		&CustomClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("%v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
	if err != nil || !Token.Valid {
		return nil, err
	}
	return Token, err
}

func Encode(req request.Encode, secret string) (string, error) {
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
	Signed, err := Token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return Signed, err
}
