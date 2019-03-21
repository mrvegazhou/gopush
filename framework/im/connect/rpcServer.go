package connect

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gopush/const"
	"gopush/framework/db/imctx"
	"gopush/framework/db"
	"gopush/framework/helper"
	"gopush/framework/im/public/transfer"
	"gopush/httpserver/models/im"
	"gopush/httpserver/service/im"
	"errors"
)
type LogicRPCer interface {
	// SignIn 设备登录
	SignIn(ctx *imctx.Context, signIn transfer.SignIn) (*transfer.SignInACK, error)
	// SyncTrigger 消息同步触发
	SyncTrigger(ctx *imctx.Context, trigger transfer.SyncTrigger) error
	// MessageSend 消息发送
	MessageSend(ctx *imctx.Context, send transfer.MessageSend) error
	// MessageACK 消息投递回执
	MessageACK(ctx *imctx.Context, ack transfer.MessageACK) error
	// OffLine 下线
	OffLine(ctx *imctx.Context, deviceId int64, userId int64) error
}

type logicRPC struct{}

var LogicRPC = new(logicRPC)

func (s *logicRPC) SignIn(ctx *imctx.Context, signIn transfer.SignIn) (*transfer.SignInACK, error) {
	device, err := imModel.DeviceDao.Get(ctx.Session, signIn.DeviceId)
	if err == gorm.ErrRecordNotFound{
		return &transfer.SignInACK{
			Code:    constdefine.IMCodeSignInFail,
			Message: "fail",
		}, nil
	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var code int
	var message string
	if device.UserId == signIn.UserId && device.Token == signIn.Token {
		err = imModel.DeviceDao.UpdateStatus(ctx.Session, signIn.DeviceId, constdefine.IMDeviceOnline)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		code = constdefine.IMCodeSignInSuccess
		message = "success"
	} else {
		code = constdefine.IMCodeSignInFail
		message = "fail"
	}

	return &transfer.SignInACK{
		Code:    code,
		Message: message,
	}, err
}

func (s *logicRPC) MessageSend(ctx *imctx.Context, send transfer.MessageSend) error {
	var err error
	send.MessageId = db.LidGenId.Get()
	fmt.Println("消息发送",
		"device_id", send.SenderDeviceId,
		"user_id", send.SenderUserId,
		"message_id", send.MessageId,
		"send_sequence", send.SendSequence)
	// 检查消息是否重复发送
	sendSequence, err := imModel.DeviceSendSequenceDao.Get(ctx.Session, send.SenderDeviceId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	ack := transfer.MessageSendACK{
		MessageId:    send.MessageId,
		DeviceId:     send.SenderDeviceId,
		SendSequence: send.SendSequence,
		Code:         constdefine.IMCCodeSuccess,
	}
	if send.SendSequence <= sendSequence.SendSequence {
		err = ConnectRPC.SendMessageSendACK(ack)
		helper.PrintErr(err)
		return nil
	}

	err = imModel.DeviceSendSequenceDao.UpdateSendSequence(ctx.Session, send.SenderDeviceId, send.SendSequence)
	helper.PrintErr(err)

	if send.ReceiverType == constdefine.IMReceiverUser {
		err = MessageServer.SendToFriend(ctx.Session, send)
	}

	if send.ReceiverType == constdefine.IMReceiverGroup {
		err = MessageServer.SendToGroup(ctx.Session, send)
	}

	if err != nil {
		return nil
	}

	// 消息发送回执
	err = ConnectRPC.SendMessageSendACK(ack)
	return nil
}


type messageServer struct{}
var MessageServer = new(messageServer)

func (*messageServer) SendToFriend(session *db.Session, send transfer.MessageSend) error {
	friend, err := imModel.FriendDao.Get(session, send.SenderUserId, send.ReceiverId)
	if err == gorm.ErrRecordNotFound || *friend == (imModel.Friend{}) {
		helper.PrintErr(err)
		return errors.New(constdefine.GetMsg(constdefine.IS_NOT_FRIEND))
	}
	if err != nil {
		helper.PrintErr(err)
		return err
	}

	selfSequence, err := imService.UserRequenceService.GetNext(session, send.SenderUserId)
	selfMessage := imModel.Message{
		MessageId:      send.MessageId,
		UserId:         send.SenderUserId,
		SenderType:     constdefine.IMSenderTypeUser,
		SenderId:       send.SenderUserId,
		SenderDeviceId: send.SenderDeviceId,
		ReceiverType:   int(send.ReceiverType),
		ReceiverId:     send.ReceiverId,
		Type:           int(send.Type),
		Content:        send.Content,
		Sequence:       selfSequence,
		SendTime:       send.SendTime,
	}

	// 发给发送者
	err = MessageServer.SendToUser(session, send.SenderUserId, &selfMessage)
	helper.PrintErr(err)

	friendSequence, err := imService.UserRequenceService.GetNext(session, send.ReceiverId)
	helper.PrintErr(err)

	friendMessage := imModel.Message{
		MessageId:      send.MessageId,
		UserId:         send.ReceiverId,
		SenderType:     constdefine.IMSenderTypeUser,
		SenderId:       send.SenderUserId,
		SenderDeviceId: send.SenderDeviceId,
		ReceiverType:   int(send.ReceiverType),
		ReceiverId:     send.ReceiverId,
		Type:           int(send.Type),
		Content:        send.Content,
		Sequence:       friendSequence,
		SendTime:       send.SendTime,
	}

	// 发给接收者
	err = MessageServer.SendToUser(session, send.ReceiverId, &friendMessage)

	return nil
}


// SendToUser 消息发送至用户
func (*messageServer) SendToUser(session *db.Session, userId int64, message *imModel.Message) error {
	_, err := imModel.MessageDao.Add(session, *message)
	helper.PrintErr(err)

	messageItem := transfer.MessageItem{
		MessageId:      message.MessageId,
		SenderType:     message.SenderType,
		SenderId:       message.SenderId,
		SenderDeviceId: message.SenderDeviceId,
		ReceiverType:   message.ReceiverType,
		ReceiverId:     message.ReceiverId,
		Type:           message.Type,
		Content:        message.Content,
		Sequence:       message.Sequence,
		SendTime:       message.SendTime,
	}

	// 查询用户在线设备
	devices, err := imModel.DeviceDao.ListOnlineByUserId(session, userId)
	helper.PrintErr(err)

	for _, v := range devices {
		message := transfer.Message{DeviceId: v.Id, Type: constdefine.IMMessageTypeMail, Messages: []transfer.MessageItem{messageItem}}
		ConnectRPC.SendMessage(message)
		fmt.Println("消息投递",
			"device_id:", message.DeviceId,
			"user_id", userId,
			"type", message.Type,
			"messages", message.GetLog())
	}
	return nil
}


func (*messageServer) SendToGroup(session *db.Session, send transfer.MessageSend) error {
	count, err := imModel.GroupUserDao.GetUserInGroup(session, send.ReceiverId, send.SenderUserId)
	helper.PrintErr(err)

	if count == 0 {
		helper.PrintErrMsg(constdefine.IM_ERROR_NOT_IN_GROUP, nil)
	}

	group, err := imService.GroupService.Get(session, send.ReceiverId)
	for _, user := range group.Users {
		sequence, err := imService.UserRequenceService.GetNext(session, user.Id)
		helper.PrintErr(err)

		message := imModel.Message{
			MessageId:      send.MessageId,
			UserId:         user.Id,
			SenderType:     constdefine.IMSenderTypeUser,
			SenderId:       send.SenderUserId,
			SenderDeviceId: send.SenderDeviceId,
			ReceiverType:   int(send.ReceiverType),
			ReceiverId:     send.ReceiverId,
			Type:           int(send.Type),
			Content:        send.Content,
			Sequence:       sequence,
			SendTime:       send.SendTime,
		}

		err = MessageServer.SendToUser(session, user.Id, &message)
		helper.PrintErr(err)
	}
	return nil
}