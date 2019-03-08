package connect

import (
	"fmt"
	"gopush/conf"
	"net"
	"runtime"
)

type TCPServer struct {
	Address      string // 端口
	MaxConnCount int    // 最大连接数
	AcceptCount  int    // 接收建立连接的groutine数量
}

func NewTCPServer(conf *conf.MainConfig) *TCPServer {
	return &TCPServer{
		Address:      conf.Tcp.Address,
		MaxConnCount: conf.Tcp.MaxConnCount,
		AcceptCount:  conf.Tcp.AcceptCount,
	}
}

func (t *TCPServer) Start() {
	addr, err := net.ResolveTCPAddr("tcp", t.Address)
	if err != nil {
		fmt.Println(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for i := 0; i < t.AcceptCount; i++ {
		go t.Accept(listener)
	}
}

func (t *TCPServer) Accept(listener *net.TCPListener) {

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = conn.SetKeepAlive(true)
		if err != nil {
			fmt.Println(err)
		}

	}

}

// RecoverPanic 恢复panic
func RecoverPanic() {
	err := recover()
	if err != nil {
		fmt.Println(err)
	}

}

// PrintStaStack 打印Panic堆栈信息
func GetPanicInfo() string {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return fmt.Sprintf("%s", buf[:n])
}