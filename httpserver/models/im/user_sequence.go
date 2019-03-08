package imModel

import (
	"gopush/framework/db"
	"time"
)

type UserSequence struct {
	Id	int64     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserId	int64	`json:"user_id"`
	Seq	int64	`json:"seq"`
	CreateTime    time.Time `json:"create_time"`    // 创建时间
	UpdateTime    time.Time `json:"update_time"`    // 更新时间
}

func (UserSequence) TableName() string {
	return "t_user_sequence"
}

type userSequenceDao struct{}

var UserSequenceDao = new(userSequenceDao)

func (*userSequenceDao) Add(session *db.Session, userId int64, seq int64) (int64, error) {
	userSequence := UserSequence{
		UserId:userId,
		Seq:seq,
	}
	if err := session.DB.Create(userSequence).Error; err != nil {
		return -1, err
	}
	return userSequence.Id, nil
}
