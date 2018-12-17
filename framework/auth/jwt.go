package jwtauth

import (
	"../../../const/code"
	"../../../const/msg"
	"../http/httputil"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
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
		payload.Username,
		payload.Password,
		payload.Id,
		time.Now().Unix(),
		time.Now().Add(exp).Unix(),
		jwt.StandardClaims{
			ExpiresAt: exp,
			Issuer:    iss,
			IssuedAt:	time.Now().Unix()
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySalt)
}

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
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
	if claims["iss"] != iss {
		return nil, constdefine.GetMsg(constdefine.ERROR_AUTH_CHECK_TOKEN_FAIL)
	}
	return claims, nil
}
