package baseController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Controller struct {
	Data interface{}
}

func ResultOk(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "data": data, "msg": ""})
}

func (this *Controller) AjaxData(ctx *gin.Context) {
	ResultOk(ctx, this.Data)
}

func (this *Controller) Redirect(ctx *gin.Context, uri string) {
	ctx.Redirect(302, uri)
}

func NoMethod(ctx *gin.Context) {
	uri := ctx.Request.RequestURI
	fmt.Printf("NoMethod" + uri)
	uri = strings.TrimLeft(uri, "/")
	uri = strings.TrimSuffix(uri, ".shtml")
	ctx.HTML(200, uri+".html", "Q")
}