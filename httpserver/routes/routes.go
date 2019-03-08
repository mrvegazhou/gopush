package routes

import (
	"github.com/gin-gonic/gin"
	"gopush/conf"
	"gopush/const"
	"gopush/framework/db"
	"gopush/framework/db/imctx"
	"gopush/framework/helper"
	"gopush/framework/http"
	"gopush/httpserver/controllers"
	"gopush/httpserver/controllers/im"
	"gopush/httpserver/service/im"
	"net/http"
	"strconv"
)

// 权限校验
func verify(c *imctx.IMContext, ctx *imctx.Context) {
	deviceIdStr := c.GetHeader("device_id")
	token := c.GetHeader("token")
	path := c.Request.URL.Path
	if path == "/device" {
		return
	}
	deviceId, err := strconv.ParseInt(deviceIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, httphelper.NewWithError(constdefine.IM_UNAUTHORIZED, constdefine.GetMsg(constdefine.IM_UNAUTHORIZED)))
		c.Abort()
		return
	}
	userId, err := imService.AuthService.Auth(ctx, deviceId, token)
	if err != nil {
		c.JSON(http.StatusOK, httphelper.NewWithError(constdefine.IM_UNAUTHORIZED, constdefine.GetMsg(constdefine.IM_UNAUTHORIZED)))
		c.Abort()
		return
	}
	c.Keys = make(map[string]interface{}, 2)
	c.Keys[constdefine.KeyDeviceId] = deviceId
	if path != "/user" && path != "/user/signin" {
		if userId == 0 {
			c.JSON(http.StatusOK, httphelper.NewWithError(constdefine.IM_DEVICE_NOT_BIND_USER, constdefine.GetMsg(constdefine.IM_DEVICE_NOT_BIND_USER)))
			c.Abort()
			return
		}
		c.Keys[constdefine.KeyUserId] = userId
	}
	c.Next()
}

func InitHandler(conf *conf.MainConfig, engine *gin.Engine) *imctx.Context {
	session := imctx.NewContext(db.ConnectDB(conf), conf)
	//验证过滤
	engine.Use(helper.Handler(verify, session))
	return session
}

func CreateRouter(router *gin.Engine, ctx *imctx.Context) {
	new(baseController.ChatController).Router(router, ctx.Conf)
	new(im.DeviceController).Router(router, ctx)
	new(im.UserController).Router(router, ctx)
	new(im.FriendController).Router(router, ctx)
}
