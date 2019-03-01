package imModel

import (
	"gopush/framework/db"
	"time"
)

type DeviceSyncSequence struct {
	Id	int64	`json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DeviceId	int64	`json:"device_id"`
	SyncSequence  int64		`json:"sync_sequence"`
	CreateTime	time.Time	`json:"create_time"`
	UpdateTime	time.Time	`json:"update_time"`
}

func (DeviceSyncSequence) TableName() string {
	return "t_device_sync_sequence"
}

type deviceSyncSequenceDao struct{}

var DeviceSyncSequenceDao = new(deviceSyncSequenceDao)

func (*deviceSyncSequenceDao) Add(session *db.Session, deviceId int64, syncSequence int64) error {
	deviceSyncSequence := DeviceSyncSequence{DeviceId:deviceId, SyncSequence:syncSequence}
	if err := session.DB.Create(&deviceSyncSequence).Error; err != nil {
		return err
	}
	return nil
}
