package routes

import (
	"../controllers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	//"../controllers/im"
	"../../framework/db/imctx"
	"../../framework/http"
	"../../const"
	"../../conf"
	"../../framework/db"
)



func handler(handler imctx.HandlerFunc, conf *conf.MainConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := new(imctx.IMContext)
		context.Context = c
		if deviceId, ok := c.Keys[constdefine.KeyDeviceId]; ok {
			context.DeviceId = deviceId.(int64)
		}
		if userId, ok := c.Keys[constdefine.KeyUserId]; ok {
			context.UserId = userId.(int64)
		}
		handler(context)
	}
}
// 权限校验
func verify(c *imctx.IMContext) {
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

	userId, err := service.AuthService.Auth(c.Session, deviceId, token)


}

func initVerify(conf *conf.MainConfig) {
	imctx.NewContext(db.ConnectDB(conf), conf)
	gin.New().Use(handler(verify, conf))
}

func CreateRouter(router *gin.Engine, conf *conf.MainConfig) {
	new(baseController.ChatController).Router(router, conf)
	//new(im.DeviceController).Router(router, conf)

}
