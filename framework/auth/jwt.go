package jwtauth

import (
	"../../const"
	jwt "github.com/dgrijalva/jwt-go"
	"time"
)

//JWT的签发者
var iss = "iss"
var exp = time.Hour * 24 * 90

type MyCustomClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Id       int    `json:"id"`
	jwt.StandardClaims
}

var mySalt = []byte("mySalt")

func GenerateToken(payload interface{}) (string, error) {
	claims := MyCustomClaims{
		Username: payload.Username,
		Password: payload.Password,
		Id: payload.Id,
		ExpiresAt: time.Now().Add(exp).Unix(),
		Issuer:    iss,
		IssuedAt:	time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySalt)
}

func ParseToken(tokenStr string) (*MyCustomClaims, string) {
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySalt, nil
	})
	if err != nil {
		return nil, constdefine.GetMsg(constdefine.ERROR_AUTH_CHECK_TOKEN_FAIL)
	}
	if jwt.SigningMethodHS256.Alg() != token.Header["alg"] {
		return nil, constdefine.GetMsg(constdefine.ERROR_AUTH_CHECK_TOKEN_FAIL)
	}
	if !token.Valid {
		return nil, constdefine.GetMsg(constdefine.ERROR_AUTH_CHECK_TOKEN_FAIL)
	}
	claims := token.Claims.(*MyCustomClaims)
	if claims.Issuer != iss {
		return nil, constdefine.GetMsg(constdefine.ERROR_AUTH_CHECK_TOKEN_FAIL)
	}
	return claims, ""
}
