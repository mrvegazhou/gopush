package imModel

import (
	"gopush/framework/db"
	"time"
)

// Friend 好友关系

type Friend struct {
	Id         int64     `json:"id"`          // 自增主键
	UserId     int64     `json:"user_id"`     // 账户id
	FriendId   int64     `json:"friend_id"`   // 好友账户id
	Label      string    `json:"label"`       // 备注，标签
	CreateTime int64 	 `json:"create_time"` // 创建时间
	UpdateTime int64 	 `json:"update_time"` // 更新时间
}

type UserFriend struct {
	UserId   int64  `json:"user_id"` // 用户id
	Label    string `json:"lable"`   // 用户对好友的标签
	Mobile   string `json:"mobile"`  // 手机号
	Nickname string `json:"nickname"`    // 昵称
}

func (Friend) TableName() string {
	return "t_friend"
}

type friendDao struct{}
var FriendDao = new(friendDao)
func (*friendDao) Add(session *db.Session, friend *Friend) (int64, error) {
	now := time.Now().UnixNano()
	friend.CreateTime = now
	friend.UpdateTime = now
	if err := session.DB.Create(friend).Error; err != nil {
		return -1, err
	}
	return friend.Id, nil
}

func (*friendDao) Get(session *db.Session, userId int64, friendId int64) (*Friend, error) {
	var friend Friend
	err := session.DB.Model(&friend).Select("id,user_id,friend_id,label,create_time,update_time").Where("user_id = ? and friend_id = ?", userId, friendId).Scan(&friend).Error
	return &friend, err
}

func (*friendDao) Delete(session *db.Session, userId, friendId int64) error {
	return session.DB.Where("user_id = ? and friend_id = ?",userId, friendId).Delete(&Friend{}).Error
}

func (*friendDao) ListUserFriend(session *db.Session, userId int64) ([]UserFriend, error) {
	rows, err := session.DB.Raw(`select f.label, u.id, u.mobile, u.nickname from t_friend f left join t_user u on f.friend_id = u.id where f.user_id = ?`, userId).Rows()
	defer rows.Close()
	friends := make([]UserFriend, 0, 5)
	for rows.Next() {
		var user UserFriend
		err := rows.Scan(&user.Label, &user.UserId, &user.Mobile, &user.Nickname)
		if err != nil {
			return nil, err
		}
		friends = append(friends, user)
	}
	return friends, err
}