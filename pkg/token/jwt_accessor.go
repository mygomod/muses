package token

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

const AccessTokenIss = "github.com/goecology/egoshop"

const AccessTokenKey = "ecologysK#xo"

type JwtTokenAccessor struct {
}

func (g *JwtTokenAccessor) EncodeAccessToken(jwtId int, uid int, startTime int64) (tokenStr string, err error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["jti"] = jwtId
	claims["iss"] = AccessTokenIss
	claims["sub"] = uid
	claims["iat"] = startTime
	claims["exp"] = startTime + AccessTokenExpireInterval
	jwtToken.Claims = claims

	tokenStr, err = jwtToken.SignedString([]byte(AccessTokenKey))
	if err != nil {
		return
	}
	return
}

func (g *JwtTokenAccessor) DecodeAccessToken(tokenStr string) (resp jwt.MapClaims, err error) {
	tokenParse, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(AccessTokenKey), nil
	})
	if err != nil {
		return
	}
	var flag bool
	resp, flag = tokenParse.Claims.(jwt.MapClaims)
	if !flag {
		err = errors.New("assert error")
		return
	}
	return
}
