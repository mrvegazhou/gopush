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
		err = service.MessageService.SendToFriend(ctx, send)
	}
}