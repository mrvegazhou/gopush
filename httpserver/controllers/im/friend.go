package im

import (
	"github.com/gin-gonic/gin"
	"gopush/framework/helper"
	"gopush/framework/db/imctx"
	"gopush/httpserver/service/im"
)

type FriendController struct {
	Controller
}

func (ctrl *FriendController) Router(router *gin.Engine, ctx *imctx.Context) {
	g := router.Group("/friend")
	g.POST("/all", helper.Handler(ctrl.All, ctx))
}

func (ctrl *FriendController) All(c *imctx.IMContext, ctx *imctx.Context) {
	data, err := imService.FriendService.ListUserFriend(ctx, c.UserId)
	ctrl.resultOk(c.Context, data, err)
}

