package imModel

import (
	"gopush/framework/db"
	"time"
)

type DeviceSendSequence struct {
	Id	int64	`json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DeviceId	int64	`json:"device_id"`
	SendSequence  int64		`json:"send_sequence"`
	CreateTime	time.Time	`json:"create_time"`
	UpdateTime	time.Time	`json:"update_time"`
}

type deviceSendSequenceDao struct{}
var DeviceSendSequenceDao = new(deviceSendSequenceDao)

func (*deviceSendSequenceDao) Add(session *db.Session, deviceId int64, sendSequence int64) error {
	deviceSendSequence := DeviceSendSequence{DeviceId:deviceId, SendSequence:sendSequence}
	if err := session.DB.Create(&deviceSendSequence).Error; err != nil {
		return err
	}
	return nil
}