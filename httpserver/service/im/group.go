package imService

import (
	"gopush/framework/db"
	"gopush/framework/helper"
	"gopush/httpserver/models/im"
)

type groupService struct{}
var GroupService = new(groupService)

//获取群组的用户信息
func (*groupService) Get(session *db.Session, id int64) (*imModel.GroupUsers, error) {
	group, err := imModel.GroupDao.Get(session, id)
	helper.PrintErr(err)

	group.Users, err = imModel.GroupUserDao.GetListGroupUser(session, id)
	helper.PrintErr(err)

	return group, err
}