package imClient

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"gopush/const"
	"gopush/framework/im/connect"
	"net"
	"gopush/framework/im/public/proto"
	"time"
)

type TcpClient struct {
	DeviceId     int64
	UserId       int64
	Token        string
	SendSequence int64
	SyncSequence int64
	codec        *connect.Codec
}

func (c *TcpClient) Start() {
	conn, err := net.Dial("tcp", "localhost:50002")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.codec = connect.NewCodec(conn)
	c.SignIn()
	go c.Receive()
}

func (c *TcpClient) SignIn() {
	signIn := pb.SignIn{
		DeviceId: c.DeviceId,
		UserId:   c.UserId,
		Token:    c.Token,
	}
	signInBytes, err := proto.Marshal(&signIn)
	if err != nil {
		fmt.Println(err)
		return
	}
	pack := connect.Package{Code: 1, Content: signInBytes}
	c.codec.Eecode(pack, 10*time.Second)
}

func (c *TcpClient) Receive() {
	for {
		_, err := c.codec.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			pack, ok := c.codec.Decode()
			if ok {
				c.HandlePackage(*pack)
				continue
			}
			break
		}

	}
}

func (c *TcpClient) HandlePackage(pack connect.Package) error {
	switch pack.Code {
		case constdefine.IMCodeSignInACK:
			ack := pb.SignInACK{}
			err := proto.Unmarshal(pack.Content, &ack)
			if err != nil {
				fmt.Println(err)
				return err
			}
			if ack.Code == 1 {
				fmt.Println("设备登录成功")
				return nil
			}
			fmt.Println("设备登录失败")

		case constdefine.IMCodeMessage:
			message := pb.Message{}
			err := proto.Unmarshal(pack.Content, &message)
			if err != nil {
				fmt.Println(err)
				return err
			}


	}
}

func (c *TcpClient) SendMessage() {

}