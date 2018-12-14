package routes

import (
	"github.com/gin-gonic/gin"
)

func createRouter(router *gin.Engine, conf interface{}) {
	new(controllers.chat).Router(router, conf)
}
