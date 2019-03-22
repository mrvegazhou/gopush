package imTest
import (
	"gopush/test/im/client"
	)
func main() {
	TestClient()
}

func TestClient() {
	client := imClient.TcpClient{}
	client.Start()
	for {
		client.SendMessage()
	}
}