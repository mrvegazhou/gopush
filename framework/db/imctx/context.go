package imctx

import (
	"gopush/framework/db"
	"gopush/conf"
	"github.com/gin-gonic/gin"
)

type Context struct {
	Session *db.Session
	Conf *conf.MainConfig
}

func NewContext(Session *db.Session, Conf *conf.MainConfig) *Context {
	return &Context{Session: Session, Conf: Conf}
}

type IMContext struct {
	*gin.Context
	DeviceId int64 // 设备id
	UserId   int64 // 用户id
}
type HandlerFunc func(*IMContext, *Context)
