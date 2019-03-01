package imModel

import (
	"time"
	"gopush/framework/db"
)

type User struct {
	Id	int64	`json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Mobile	string	`json:"mobile"`
	Nickname	string	`json:"nickname"`
	Password	string `json:"password"`
	CreateTime    time.Time `json:"create_time"`    // 创建时间
	UpdateTime    time.Time `json:"update_time"`    // 更新时间
}

// SignIn 登录结构体
type SignIn struct {
	Number   string `json:"number"`
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

func (*userDao) Add(session *db.Session, user User) (int64, error) {
	if err := session.DB.Create(user).Error; err != nil {
		return -1, err
	}
	return user.Id, nil
}