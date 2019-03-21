package imModel

import (
	"gopush/framework/db"
	"gopush/framework/helper"
)

// Group 群组
type Group struct {
	Id        int64       `json:"id"`    // 群组id
	Name      string      `json:"name"`  // 组名
	CreateTime int64 	 `json:"create_time"` // 创建时间
	UpdateTime int64 	 `json:"update_time"` // 更新时间
}

type GroupUsers struct {
	Id        int64       `json:"id"`    // 群组id
	Name      string      `json:"name"`  // 组名
	Users []User `json:"users"` // 群组用户
	CreateTime int64 	 `json:"create_time"` // 创建时间
	UpdateTime int64 	 `json:"update_time"` // 更新时间
}

type GroupUserUpdate struct {
	GroupId int64   `json:"group_id"` // 群组名称
	UserIds []int64 `json:"user_ids"` // 群组成员
}

type groupDao struct{}
var GroupDao = new(groupDao)

func (*groupDao) Get(session *db.Session, id int64) (*GroupUsers, error) {
	var groupUsers GroupUsers
	err := session.DB.Model(Group{}).Select("id,name,create_time,update_time").Where("id = ?", id).Scan(&groupUsers).Error
	return &groupUsers, err
}

func (*groupDao) Add(session *db.Session, name string) (int64, error) {
	group := new(Group)
	now := helper.NowUnixTime()
	group.CreateTime = now
	group.UpdateTime = now
	group.Name = name
	if err := session.DB.Create(group).Error; err != nil {
		return -1, err
	}
	return group.Id, nil
}

