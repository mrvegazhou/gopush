package imModel

import (
	"gopush/framework/db"
	"gopush/framework/helper"
)

type DeviceSendSequence struct {
	Id	int64	`json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DeviceId	int64	`json:"device_id"`
	SendSequence  int64		`json:"send_sequence"`
	CreateTime	int64	`json:"create_time"`
	UpdateTime	int64	`json:"update_time"`
}

func (DeviceSendSequence) TableName() string {
	return "t_device_send_sequence"
}

type deviceSendSequenceDao struct{}
var DeviceSendSequenceDao = new(deviceSendSequenceDao)

func (*deviceSendSequenceDao) Add(session *db.Session, deviceId int64, sendSequence int64) error {
	now := helper.NowUnixTime()
	deviceSendSequence := DeviceSendSequence{DeviceId:deviceId, SendSequence:sendSequence, CreateTime:now, UpdateTime:now}
	if err := session.DB.Create(&deviceSendSequence).Error; err != nil {
		return err
	}
	return nil
}

func (*deviceSendSequenceDao) UpdateSendSequence(session *db.Session, deviceId int64, sendSequence int64) error {
	var deviceSendSequence DeviceSendSequence
	update := helper.NowUnixTime()
	if err := session.DB.Model(deviceSendSequence).Where("send_sequence = ?", sendSequence).Updates(map[string]interface{}{"device_id": deviceId, "update_time": update}).Error; err != nil {
		return err
	}
	return nil
}

func (*deviceSendSequenceDao) Get(session *db.Session, id int64) (*DeviceSendSequence, error) {
	var deviceSendSequence DeviceSendSequence
	err := session.DB.Model(&deviceSendSequence).Select("send_sequence from t_device_send_sequence").Where("device_id=?", id).Scan(&deviceSendSequence).Error
	return &deviceSendSequence, err
}