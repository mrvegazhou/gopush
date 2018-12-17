package routes

import (
	"../controllers"
	"github.com/gin-gonic/gin"
)

func CreateRouter(router *gin.Engine, conf interface{}) {
	new(baseController.ChatController).Router(router, conf)
}
