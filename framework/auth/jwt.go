package jwtauth

import (
	"../http/httputil"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	Id       int    `json:"id"`
	jwt.StandardClaims
}

func GenerateToken(signingKey string, payload interface{}) (string, error) {
	claims := MyCustomClaims{
		payload.Username,
		payload.Id,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func JwtMiddleware(next http.Handler, conf interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		fmt.Printf("err:", err)
		// 	}
		// }()
		tokenStr := r.Header.Get("authorization")
		if tokenStr == "" {
			httphelper.ResponseWithJson(w, http.StatusUnauthorized, httphelper.Response{Code: http.StatusUnauthorized, Msg: "not authorized"})
		} else {
			token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					httphelper.ResponseWithJson(w, http.StatusUnauthorized, httphelper.Response{Code: http.StatusUnauthorized, Msg: "not authorized"})
					return nil, fmt.Errorf("not authorization")
				}
				return []byte(conf.Jwt.Key), nil
			})
			var msg
			if token.Valid {
				msg := "You look nice today"
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					msg := "That's not even a token"
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					msg := "Timing is everything"
				} else {
					msg := fmt.Sprintf("%s:%s", "Couldn't handle this token:" , err)
				}
			} else {
				msg := fmt.Sprintf("%s:%s", "Couldn't handle this token:", err)
			}
			return nil, fmt.Errorf(msg)
		}
	})
}
