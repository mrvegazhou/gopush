package imModel

import (
	"time"
	"gopush/framework/db/imctx"
)

type UserSequence struct {
	Id	int64     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	UserId	int64	`json:"user_id"`
	Sequence	int64	`json:"sequence"`
	CreateTime    time.Time `json:"create_time"`    // 创建时间
	UpdateTime    time.Time `json:"update_time"`    // 更新时间
}

type userSequenceDao struct{}

var UserSequenceDao = new(userSequenceDao)

func (*userSequenceDao) Add(ctx *imctx.Context, userId int64, sequence int64) (int64, error) {
	if err := session.DB.Create(device).Error; err != nil {
		return -1, err
	}
	return device.Id, nil
}
