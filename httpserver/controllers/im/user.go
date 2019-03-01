package im

import (
	"github.com/gin-gonic/gin"
	"gopush/framework/helper"
	"gopush/framework/db/imctx"
)

type UserController struct {
	Controller
}

func (ctrl *UserController) Router(router *gin.Engine, ctx *imctx.Context) {
	router.POST("/user", helper.Handler(ctrl.Regist, ctx))
}

func (ctrl *UserController) Regist(c *imctx.IMContext, ctx *imctx.Context) {

}