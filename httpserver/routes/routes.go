package routes

import (
	"../controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"../../framework/db/imctx"
	"../../framework/http"
	"../../const"
	"../../conf"
	"../../framework/db"
	"../service"
)

func handler(handler imctx.HandlerFunc, session *imctx.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(imctx.IMContext)
		context.Context = c
		if deviceId, ok := c.Keys[constdefine.KeyDeviceId]; ok {
			context.DeviceId = deviceId.(int64)
		}
		if userId, ok := c.Keys[constdefine.KeyUserId]; ok {
			context.UserId = userId.(int64)
		}
		handler(context, session)
	}
}
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

	userId, err := service.AuthService.Auth(ctx, deviceId, token)
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

func initVerify(conf *conf.MainConfig) {
	gin.New().Use(handler(verify, imctx.NewContext(db.ConnectDB(conf), conf)))
}

func CreateRouter(router *gin.Engine, conf *conf.MainConfig) {
	new(baseController.ChatController).Router(router, conf)
	//new(im.DeviceController).Router(router, conf)

}
