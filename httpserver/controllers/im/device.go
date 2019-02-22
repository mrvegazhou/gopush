package im

import (
	"github.com/gin-gonic/gin"
	"gopush/conf"
)

type DeviceController struct {
	Controller
}

func (ctrl *DeviceController) Router(router *gin.Engine, conf *conf.MainConfig) {

	router.POST("/device", handler(ctrl.Regist))
}

func (DeviceController) Regist(c *context) {

}



