package baseController

import (
	"../../framework/pubsub/chat/longpolling"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ChatController struct {
	Controller
}

// func (ctrl *ChatController) before() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		uri := ctx.Request.RequestURI
// 		fmt.Print(uri)
// 		if 1 == 1 {
// 			ctx.Next()
// 		}
// 		return
// 	}
// }

func (ctrl *ChatController) Router(router *gin.Engine, conf interface{}) {
	router.GET("/lp", ctrl.showIndex)
	router.GET("/lp/fetch", ctrl.fetch)
	router.POST("/lp/post", ctrl.post)
}

func (ctrl *ChatController) fetch(ctx *gin.Context) {
	lastReceived := ctx.DefaultQuery("lastReceived", "111") //time.Now().Unix()
	fmt.Printf("lastReceived:", lastReceived)
	lastReceivedInt, _ := strconv.Atoi(lastReceived)
	ctrl.Data = longpolling.FetchMsgs(lastReceivedInt)
	ctrl.AjaxData(ctx)
}

func (ctrl *ChatController) showIndex(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.HTML(200, "longpolling.html", gin.H{})
}

func (ctrl *ChatController) post(ctx *gin.Context) {
	uname := ctx.PostForm("uname")
	content := ctx.PostForm("content")
	longpolling.PostMsg(uname, content)
}

func (ctrl *ChatController) Redirect(ctx *gin.Context) {
	ctx.Redirect(302, "/")
}
