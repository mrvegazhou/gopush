package connect

import (
	"fmt"
	"gopush/const"
	"io"
	"net"
	"strings"
	"time"
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

func (c *ConnContext) DoConn() {
	defer RecoverPanic()
	for {
		err := c.Codec.Conn.SetReadDeadline(time.Now().Add(constdefine.IMReadDeadline))

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