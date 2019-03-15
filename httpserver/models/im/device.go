package imModel

import (
	"gopush/framework/db"
	"gopush/framework/helper"
)

type Device struct {
	Id            int64     `json:"id" gorm:"primary_key;AUTO_INCREMENT"`             // 设备id
	UserId        int64     `json:"user_id"`        // 用户id
	Token         string    `json:"token"`          // 设备登录的token
	Type          int       `json:"type"`           // 设备类型,1:Android；2：IOS；3：Windows; 4：MacOS；5：Web
	Brand         string    `json:"brand"`          // 手机厂商
	Model         string    `json:"model"`          // 机型
	SystemVersion string    `json:"system_version"` // 系统版本
	APPVersion    string    `json:"app_version"`    // APP版本
	Status        int       `json:"state"`          // 在线状态，0：不在线；1：在线
	CreateTime    int64 	 `json:"create_time"`    // 创建时间
	UpdateTime    int64 	 `json:"update_time"`    // 更新时间
}

func (Device) TableName() string {
	return "t_device"
}

type deviceDao struct{}
var DeviceDao = new(deviceDao)

// Get 获取设备
func (*deviceDao) Get(session *db.Session, id int64) (*Device, error) {
	var device Device
	//db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
	err := session.DB.Model(&device).Select("user_id,token,type,brand,model,system_version,app_version,status,create_time,update_time").Where("id=?", id).Scan(&device).Error
	return &device, err
}

func (*deviceDao) Add(session *db.Session, device *Device) (int64, error) {
	now := helper.NowUnixTime()
	device.CreateTime = now
	device.UpdateTime = now
	if err := session.DB.Create(device).Error; err != nil {
		return -1, err
	}
	return device.Id, nil
}

func (*deviceDao) UpdateUserId(session *db.Session, id, userId int64) error {
	var device Device
	update := helper.NowUnixTime()
	if err := session.DB.Model(device).Where("id = ?", id).Updates(map[string]interface{}{"user_id": userId, "update_time": update}).Error; err != nil {
		return err
	}
	return nil
}

func (*deviceDao) UpdateStatus(session *db.Session, id int64, status int) error {
	var device Device
	update := helper.NowUnixTime()
	if err := session.DB.Model(device).Where("id = ?", id).Updates(map[string]interface{}{"status": status, "update_time": update}).Error; err != nil {
		return err
	}
	return nil
}

