package connect

import (
	"github.com/golang/protobuf/proto"
	"gopush/const"
	"gopush/framework/helper"
	"gopush/framework/im/public/proto"
	"gopush/framework/im/public/transfer"
)

// ConnectRPCer 连接层接口
type ConnectRPCer interface {
	SendMessage(message transfer.Message) error
	SendMessageSendACK(ack transfer.MessageSendACK) error
}

type connectRPC struct{}
var ConnectRPC = new(connectRPC)

func (*connectRPC) SendMessageSendACK(ack transfer.MessageSendACK) error {
	content, err := proto.Marshal(&pb.MessageSendACK{SendSequence: ack.SendSequence, Code: int32(ack.Code)})
	helper.PrintErr(err)

	ctx, err := load(ack.DeviceId)
	helper.PrintErr(err)

	err = ctx.Codec.Eecode(Package{Code: constdefine.IMCodeMessageSendACK, Content: content}, constdefine.IMWriteDeadline)
	helper.PrintErr(err)

	return nil
}

// SendMessage 处理消息投递
func (*connectRPC) SendMessage(message transfer.Message) error {
	ctx, err := load(message.DeviceId)
	helper.PrintErr(err)

	messages := make([]*pb.MessageItem, 0, len(message.Messages))
	for _, v := range message.Messages {
		item := new(pb.MessageItem)
		item.MessageId = v.MessageId
		item.SenderType = int32(v.SenderType)
		item.SenderId = v.SenderId
		item.SenderDeviceId = v.SenderDeviceId
		item.ReceiverType = int32(v.ReceiverType)
		item.ReceiverId = v.ReceiverId
		item.Type = int32(v.Type)
		item.Content = v.Content
		item.SyncSequence = v.Sequence
		item.SendTime = v.SendTime / 1000000

		messages = append(messages, item)
	}
	content, err := proto.Marshal(&pb.Message{Type: message.Type, Messages: messages})
	helper.PrintErr(err)

	err = ctx.Codec.Eecode(Package{Code: constdefine.IMCodeMessage, Content: content}, constdefine.IMWriteDeadline)
	helper.PrintErr(err)

	return nil
}