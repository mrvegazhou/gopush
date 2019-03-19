package imService

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"gopush/const"
	"gopush/framework/db"
	"gopush/framework/helper"
	"gopush/framework/im/public/transfer"
	"gopush/httpserver/models/im"
)

type messageService struct{}
var MessageService = new(messageService)

func (*messageService) SendToFriend(session *db.Session, send transfer.MessageSend) error {
	friend, err := imModel.FriendDao.Get(session, send.SenderUserId, send.ReceiverId)
	if err == gorm.ErrRecordNotFound {
		helper.PrintErr(err)
		return errors.New(constdefine.GetMsg(constdefine.IS_NOT_FRIEND))
	}
	if err != nil {
		helper.PrintErr(err)
		return err
	}

}
