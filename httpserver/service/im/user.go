package imService
import (
	"database/sql"
	"fmt"
	"errors"
	"gopush/const"
	"gopush/httpserver/models/im"
	"gopush/framework/db/imctx"
)
type userService struct{}

var UserService = new(userService)

func (*userService) Regist(ctx *imctx.Context, deviceId int64, regist imModel.User) (*imModel.SignInResp, error) {
	db := ctx.Session.DB
	tx := db.Begin()
	user := imModel.User{
		Mobile:   regist.Mobile,
		Nickname: regist.Nickname,
		Password: regist.Password,
	}
	userId, err := imModel.UserDao.Add(ctx, user)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return nil, err
	}
	if userId == -1 {
		tx.Rollback()
		return nil, errors.New(constdefine.GetMsg(constdefine.IM_ERROR_USER_REGIST))
	}

	err = imModel.DeviceDao.UpdateUserId(ctx.Session, deviceId, userId)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return nil, err
	}

	err = imModel.DeviceSendSequenceDao.UpdateSendSequence(ctx.Session, deviceId, 0)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return nil, err
	}

	imModel.DeviceSyncSequenceDao.UpdateSyncSequence(ctx.Session, deviceId, 0)
	if err != nil {
		tx.Rollback()
		fmt.Println(err)
		return nil, err
	}

	tx.Commit()

	return &imModel.SignInResp{
		SendSequence: 0,
		SyncSequence: 0,
	}, nil
}

func (*userService) SignIn(ctx *imctx.Context, deviceId int64, mobile string, password string) (*imModel.SignInResp, error) {
	db := ctx.Session.DB
	tx := db.Begin()
	user, err := imModel.UserDao.GetByNumber(ctx.Session, mobile)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(constdefine.GetMsg(constdefine.ERROR_MOBILE))
		}
		return nil, err
	}
	if password != user.Password {
		return nil, errors.New(constdefine.GetMsg(constdefine.ERROR_PASSWORD))
	}

	err = imModel.DeviceDao.UpdateUserId(ctx.Session, deviceId, user.Id)
	if err != nil {
		return nil, err
	}

	err = imModel.DeviceSendSequenceDao.UpdateSendSequence(ctx.Session, deviceId, 0)
	if err != nil {
		return nil, err
	}

	maxSyncSequence, err := imModel.DeviceSyncSequenceDao.GetMaxSyncSequenceByUserId(ctx.Session, user.Id)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return &imModel.SignInResp{
		SendSequence: 0,
		SyncSequence: maxSyncSequence,
	}, nil
}
