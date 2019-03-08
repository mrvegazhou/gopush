package imService

import (
	"gopush/framework/db/imctx"
	"gopush/httpserver/models/im"
)
type friendService struct{}

var FriendService = new(friendService)

func (*friendService) ListUserFriend(ctx *imctx.Context, userId int64) ([]imModel.UserFriend, error) {
	friends, err := imModel.FriendDao.ListUserFriend(ctx.Session, userId)
	if err != nil {
		return nil, err
	}
	return friends, err
}
