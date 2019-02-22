package middlewareJwt

import (
	"../../../const"
	"../../../framework/auth"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		path := c.Request.URL.Path
		if path == "/device" {
			return
		}
		token := c.GetHeader("token")
		if token == "" {
			code = constdefine.INVALID_PARAMS
		} else {
			claims, err := jwtauth.ParseToken(token)
			if err != "" {
				code = constdefine.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = constdefine.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != constdefine.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"Code": code,
				"Msg":  constdefine.GetMsg(code),
				"Data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
