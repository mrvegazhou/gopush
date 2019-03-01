package im

import (
	"gopush/const"
	"gopush/framework/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Data interface{}
}


func (c *Controller) badParam(ctx *gin.Context, err int) {
	ctx.JSON(http.StatusOK, httphelper.NewWithError(err, constdefine.GetMsg(err)))
}

func (c *Controller) resultOk(ctx *gin.Context, data interface{}, err error) {
	if err != nil {
		ctx.JSON(http.StatusOK, httphelper.NewWithError(constdefine.IM_UNKNOWN_ERROR, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, httphelper.NewSuccess(data))
}

func (c *Controller) bindJson(ctx *gin.Context, value interface{}) error {
	err := ctx.ShouldBindJSON(value)
	if err != nil {
		ctx.JSON(http.StatusOK, httphelper.NewWithError(constdefine.IM_BAD_REQUEST, err.Error()))
		return err
	}
	return nil
}