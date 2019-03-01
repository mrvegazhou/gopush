package imService
import (
	"github.com/pkg/errors"
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
	userId, err := imModel.UserDao.Add(ctx.Session, user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if userId == -1 {
		tx.Rollback()
		return nil, errors.New(constdefine.GetMsg(constdefine.IM_ERROR_USER_REGIST))
	}

	tx.Commit()
}
