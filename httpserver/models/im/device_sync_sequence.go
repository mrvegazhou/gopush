package imModel

import (
	"gopush/framework/db"
)

type DeviceSyncSequence struct {
	Id	int64	`json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	DeviceId	int64	`json:"device_id"`
	SyncSequence  int64		`json:"sync_sequence"`
	CreateTime	int64	`json:"create_time"`
	UpdateTime	int64	`json:"update_time"`
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

func (*deviceSyncSequenceDao) UpdateSyncSequence(session *db.Session, deviceId int64, sendSequence int64) error {
	db := session.DB
	var deviceSyncSequence DeviceSyncSequence
	if err := db.Model(deviceSyncSequence).Where("sync_sequence = ?", sendSequence).Update("device_id", deviceId).Error; err != nil {
		return err
	}
	return nil
}

type MaxStruct struct {
	Max int64 `json:"max"`
}
func (*deviceSyncSequenceDao) GetMaxSyncSequenceByUserId(session *db.Session, userId int64) (int64, error) {
	var result MaxStruct
	err := session.DB.Raw(`select max(s.sync_sequence)
									from t_device d
									left join t_device_sync_sequence s on d.id = s.device_id
									where user_id = ?`, userId).Scan(&result).Error
	if err != nil {
		return -1, err
	}
	return result.Max, nil
}
