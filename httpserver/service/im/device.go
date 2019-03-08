package imService

import (
	"fmt"
	"gopush/httpserver/models/im"
	"gopush/framework/db/imctx"
	"github.com/satori/go.uuid"
	//"gopush/framework/log"
)

type deviceService struct{}
var DeviceService = new(deviceService)

// Regist 注册设备
func (*deviceService) Regist(ctx *imctx.Context, device *imModel.Device) (int64, string, error) {
	db := ctx.Session.DB
	tx := db.Begin()
	UUID, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		tx.Rollback()
		return 0, "", err
	}
	device.Token = UUID.String()

	id, err := imModel.DeviceDao.Add(ctx.Session, device)
	if err != nil {
		//logger.Sugar.Error(err)
		fmt.Printf(err.Error())
		tx.Rollback()
		return 0, "", err
	}

	err = imModel.DeviceSendSequenceDao.Add(ctx.Session, id, 0)
	if err != nil {
		//logger.Sugar.Error(err)
		tx.Rollback()
		return 0, "", err
	}

	err = imModel.DeviceSyncSequenceDao.Add(ctx.Session, id, 0)
	if err != nil {
		//logger.Sugar.Error(err)
		tx.Rollback()
		return 0, "", err
	}

	tx.Commit()
	return id, device.Token, nil
}

