package imModel

import (
	"github.com/jinzhu/gorm"
	"gopush/framework/db"
	"gopush/framework/helper"
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

func (*userSequenceDao) GetSequence(session *db.Session, userId int64) (int64, error) {
	var sequence int64
	var userSequence UserSequence
	if err := session.DB.Model(&userSequence).Select("id,user_id,friend_id,seq,create_time,update_time").Where("user_id = ?", userId).Scan(&userSequence).Error; err!=nil {
		return -1, err
	}
	return sequence, nil
}

func (*userSequenceDao) Increase(session *db.Session, userId int64) error {
	var userSequence UserSequence
	update := helper.NowUnixTime()
	if err := session.DB.Model(&userSequence).Where("user_id = ?", userId).Update(map[string]interface{}{"sequence": gorm.Expr("sequence + ?", 1), "update_time": update}).Error; err != nil {
		return err
	}
	return nil
}