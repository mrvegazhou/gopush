package imService

import (
	"database/sql"
	"errors"
	"gopush/const"
	"gopush/framework/db/imctx"
	"gopush/httpserver/models/im"
)
type authService struct{}

var AuthService = new(authService)

func (*authService) Auth(ctx *imctx.Context, deviceId int64, token string) (int64, error) {
	device, err := imModel.DeviceDao.Get(ctx.Session, deviceId)
	if err == sql.ErrNoRows {
		return 0, errors.New(constdefine.GetMsg(constdefine.IM_UNAUTHORIZED))
	}
	if err != nil {
		return 0, err
	}
	if token != device.Token {
		return 0, errors.New(constdefine.GetMsg(constdefine.IM_UNAUTHORIZED))
	}
	return device.UserId, nil
}
