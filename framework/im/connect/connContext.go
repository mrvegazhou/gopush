package connect

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopush/const"
	"gopush/framework/helper"
	"io"
	"net"
	"strings"
	"time"
	"gopush/framework/im/public/proto"
	"gopush/framework/im/public/transfer"
	"gopush/framework/db/imctx"
)

// Package 消息包
type Package struct {
	Code    int    // 消息类型
	Content []byte // 消息体
}

// ConnContext 连接上下文
type ConnContext struct {
	Codec    *Codec // 编解码器
	IsSignIn bool   // 是否登录
	DeviceId int64  // 设备id
	UserId   int64  // 用户id
}

func NewConnContext(conn *net.TCPConn) *ConnContext {
	codec := NewCodec(conn)
	return &ConnContext{Codec: codec}
}

func (c *ConnContext) DoConn(ctx *imctx.Context) {
	defer RecoverPanic()
	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(constdefine.IMReadDeadline))
		if err != nil {
			c.HandleReadErr(err)
			return
		}
		_, err = c.Codec.Read()
		if err != nil {
			c.HandleReadErr(err)
			return
		}

		for {
			message, ok := c.Codec.Decode()
			if ok {
				c.HandlePackage(message, ctx)
				continue
			}
			break
		}
	}
}

func (c *ConnContext) HandlePackage(pack *Package, ctx *imctx.Context) {
	// 未登录拦截
	if pack.Code != constdefine.IMCodeSignIn && c.IsSignIn == false {
		c.Release()
		return
	}
	switch pack.Code {
	case constdefine.IMCodeSignIn:
		c.HandlePackageSignIn(pack, ctx)
	case constdefine.IMCodeSyncTrigger:
	case constdefine.IMCodeHeadbeat:
		c.HandlePackageHeadbeat()
	case constdefine.IMCodeMessageSend:
		c.HandlePackageMessageSend(pack)
	case constdefine.IMCodeMessageACK:
	}
}

// HandlePackageSignIn 处理登录消息包
func (c *ConnContext) HandlePackageSignIn(pack *Package, ctx *imctx.Context) {
	var signIn pb.SignIn
	err := proto.Unmarshal(pack.Content, &signIn)
	if err != nil {
		fmt.Println(err)
		c.Release()
		return
	}
	transferSignIn := transfer.SignIn{
		DeviceId: signIn.DeviceId,
		UserId:   signIn.UserId,
		Token:    signIn.Token,
	}

	// 处理设备登录逻辑
	ack, err := LogicRPC.SignIn(ctx, transferSignIn)
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := proto.Marshal(&pb.SignInACK{Code: int32(ack.Code), Message: ack.Message})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = c.Codec.Eecode(Package{Code: constdefine.IMCodeSignInACK, Content: content}, constdefine.IMWriteDeadline)
	if err != nil {
		fmt.Println(err)
		return
	}

	if ack.Code == constdefine.IMCodeSignInSuccess{
		c.IsSignIn = true
		c.DeviceId = signIn.DeviceId
		c.UserId = signIn.UserId
		store(c.DeviceId, c)
	}
}

func (c *ConnContext) HandlePackageHeadbeat() {
	err := c.Codec.Eecode(Package{Code: constdefine.IMCodeHeadbeatACK, Content: []byte{}}, constdefine.IMWriteDeadline)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("心跳：", "device_id", c.DeviceId, "user_id", c.UserId)
}

func (c *ConnContext) HandlePackageMessageSend(pack *Package) {
	var send pb.MessageSend
	err := proto.Unmarshal(pack.Content, &send)
	if err != nil {
		fmt.Println(err)
		c.Release()
		return
	}
	transferSend := transfer.MessageSend{
		SenderDeviceId: c.DeviceId,
		SenderUserId:   c.UserId,
		ReceiverType:   send.ReceiverType,
		ReceiverId:     send.ReceiverId,
		Type:           send.Type,
		Content:        send.Content,
		SendSequence:   send.SendSequence,
		SendTime:		helper.UnunixTime(send.SendTime),
	}


}

// HandleReadErr 读取conn错误
func (c *ConnContext) HandleReadErr(err error) {
	fmt.Println("连接读取异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
	str := err.Error()
	// 服务器主动关闭连接
	if strings.HasSuffix(str, "use of closed network connection") {
		return
	}

	c.Release()
	// 客户端主动关闭连接或者异常程序退出
	if err == io.EOF {
		return
	}
	// SetReadDeadline 之后，超时返回的错误
	if strings.HasSuffix(str, "i/o timeout") {
		return
	}
	fmt.Println("连接读取未知异常：", "device_id", c.DeviceId, "user_id", c.UserId, "err_msg", err)
}

// Release 释放TCP连接
func (c *ConnContext) Release() {
	delete(c.DeviceId)
	err := c.Codec.Conn.Close()
	if err != nil {
		fmt.Println(err)
	}
	//err = LogicRPC.OffLine(Context(), c.DeviceId, c.UserId)
	//if err != nil {
	//	fmt.Println(err)
	//}
}