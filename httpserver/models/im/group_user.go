package imModel

import (
	"gopush/framework/db"
)

// GroupUser 群组成员
type GroupUser struct {
	UserId int64  `json:"user_id"` // 用户id
	GroupId int64  `json:"group_id"`
	Label  string `json:"label"`   // 用户标签
	CreateTime int64 	 `json:"create_time"` // 创建时间
	UpdateTime int64 	 `json:"update_time"` // 更新时间
}

func (GroupUser) TableName() string {
	return "t_group_user"
}


type groupUserDao struct{}
var GroupUserDao = new(groupUserDao)

// UserInGroup 用户是否在群组中
func (*groupUserDao) GetUserInGroup(session *db.Session, groupId int64, userId int64) (int, error) {
	var count int
	err := session.DB.Model(&GroupUser{}).Where(&GroupUser{UserId: userId, GroupId: groupId}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (*groupUserDao) GetListGroupUser(session *db.Session, id int64) ([]User, error) {
	rows, err := session.DB.Raw(`select user_id,group_id,label,create_time,update_time from t_group_user where group_id=?`, id).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	userIds := make([]int64, 0, 5)
	for rows.Next() {
		var groupUser GroupUser
		err := rows.Scan(&groupUser.UserId, &groupUser.GroupId, &groupUser.Label, &groupUser.CreateTime, &groupUser.UpdateTime)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, groupUser.UserId)
	}
	users := []User{}
	err = session.DB.Table("t_user").Where("id in (?)", userIds).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}