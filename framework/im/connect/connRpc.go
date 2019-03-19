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