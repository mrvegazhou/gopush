package connect

import (
	"fmt"
	"gopush/framework/db/imctx"
	"net"
	"runtime"
)

type TCPServer struct {
	Address      string // 端口
	MaxConnCount int    // 最大连接数
	AcceptCount  int    // 接收建立连接的groutine数量
}

func NewTCPServer(ctx *imctx.Context) *TCPServer {
	return &TCPServer{
		Address:      	ctx.Conf.Tcp.Address,
		MaxConnCount: 	ctx.Conf.Tcp.MaxConnCount,
		AcceptCount:  	ctx.Conf.Tcp.AcceptCount,
	}
}

func (t *TCPServer) Start(ctx *imctx.Context) {
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
		go t.Accept(listener, ctx)
	}

	select {}
}

func (t *TCPServer) Accept(listener *net.TCPListener, ctx *imctx.Context) {

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

		connContext := NewConnContext(conn)
		go connContext.DoConn(ctx)
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