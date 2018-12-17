package baseController

import (
	"../../../const/code"
	"../../../const/msg"
	"fmt"
	"github.com/gin-gonic/gin"
)

type JwtController struct {
	Controller
}

func (ctrl *JwtController) Router(router *gin.Engine, conf interface{}) {
	router.GET("/auth", ctrl.getAuth)
}

type auth struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

func getAuth(c *gin.Context) {
	err := c.ShouldBind(&auth)
	result := gin.H{}
	if err != nil {
		result["code"] = constdefine.InvalidParams
		result["message"] = constdefine.GetMsg(constdefine.InvalidParams)
		c.JSON(200, result)
		return
	}
}
