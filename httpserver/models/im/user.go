package imModel

import (
	"gopush/framework/db"
	"gopush/framework/db/imctx"
	"gopush/framework/helper"
)

type User struct {
	Id	int64	`json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Mobile	string	`json:"mobile"`
	Nickname	string	`json:"nickname"`
	Password	string `json:"password"`
	CreateTime    int64 `json:"create_time"`    // 创建时间
	UpdateTime    int64 `json:"update_time"`    // 更新时间
}

// SignIn 登录结构体
type SignIn struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

// SignInResp 登录响应
type SignInResp struct {
	SendSequence int64 `json:"send_sequence"` // 发送序列号
	SyncSequence int64 `json:"sync_sequence"` // 同步序列号
}

func (User) TableName() string {
	return "t_user"
}

type userDao struct{}
var UserDao = new(userDao)

func (*userDao) Add(ctx *imctx.Context, user User) (int64, error) {
	now := helper.NowUnixTime()
	user.CreateTime = now
	user.UpdateTime = now
	db := ctx.Session.DB
	db.LogMode(ctx.Conf.Postgresql.DbDebug)
	if err := db.Create(&user).Error; err != nil {
		return -1, err
	}
	return user.Id, nil
}

func (*userDao) GetByNumber(session *db.Session, mobile string) (*User, error) {
	var user User
	err := session.DB.Select("id,mobile,nickname,password").Where("mobile=?", mobile).Scan(&user).Error
	return &user, err
}