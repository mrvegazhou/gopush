package im

import (
	"github.com/gin-gonic/gin"
	"gopush/framework/db/imctx"
	"gopush/framework/helper"
	"gopush/httpserver/models/im"
	"gopush/httpserver/service/im"
)

type UserController struct {
	Controller
}

func (ctrl *UserController) Router(router *gin.Engine, ctx *imctx.Context) {
	router.POST("/user", helper.Handler(ctrl.Regist, ctx))
}

func (ctrl *UserController) Regist(c *imctx.IMContext, ctx *imctx.Context) {
	var regist imModel.User
	if ctrl.bindJson(c.Context, &regist) != nil {
		return
	}
	data, err := imService.UserService.Regist(ctx, c.DeviceId, regist)

	ctrl.resultOk(c.Context, data, err)
}

func (ctrl *UserController) SignIn(c *imctx.IMContext, ctx *imctx.Context) {
	var data struct {
		Mobile   string `json:"mobile"`
		Password string `json:"password"`
	}
	if ctrl.bindJson(c.Context, &data) != nil {
		return
	}
	var sign *imModel.SignInResp
	sign, err := imService.UserService.SignIn(ctx, c.DeviceId, data.Mobile, data.Password)
	ctrl.resultOk(c.Context, sign, err)
}