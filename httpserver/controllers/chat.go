package controllers

import (
	"./baseController"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatController struct {
	restgo.Controller
}

func (ctrl *ChatController) before() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uri := ctx.Request.RequestURI
		fmt.Print(uri)
		if 1 == 1 {
			ctx.Next()
		}
		return
	}
}

func (ctrl *ChatController) create(ctx *gin.Context) {
	ctrl.Data = []int{1, 2, 3}
	ctrl.AjaxData(ctx)
}
