package im

import (
	"github.com/gin-gonic/gin"
	"gopush/const"
	"gopush/framework/db/imctx"
	"gopush/framework/helper"
	"gopush/httpserver/service/im"
	"gopush/httpserver/models/im"
)

type DeviceController struct {
	Controller
}

func (ctrl *DeviceController) Router(router *gin.Engine, ctx *imctx.Context) {
	router.POST("/device", helper.Handler(ctrl.Regist, ctx))
}

func (ctrl *DeviceController) Regist(c *imctx.IMContext, ctx *imctx.Context) {
	var device imModel.Device
	if c.ShouldBindJSON(&device) != nil {
		return
	}

	if device.Type == 0 || device.Brand == "" || device.Model == "" || device.SystemVersion == "" || device.APPVersion == "" {
		ctrl.badParam(c.Context, constdefine.IM_BAD_REQUEST)
		return
	}

	id, token, err := imService.DeviceService.Regist(ctx, &device)
	ctrl.resultOk(c.Context, map[string]interface{}{"id": id, "token": token}, err)
}



