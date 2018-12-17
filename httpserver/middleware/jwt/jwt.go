package middlewareJwt

import (
	"../../../const/code"
	"../../../const/msg"
	"../../../framework/auth/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		token := c.Query("token")
		if token == "" {
			code = constdefine.INVALID_PARAMS
		} else {
			claims, err := jwtauth.ParseToken(token)
			if err != nil {
				code = constdefine.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = constdefine.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != constdefine.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  constdefine.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
