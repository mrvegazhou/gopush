package websocket

import (
	"../chatroom"
	"container/list"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Join(w http.ResponseWriter, r *http.Request) {
	// 连接的指针
	ws, err := upgrader.Upgrade(w, r, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(this.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}
	Join(uname, ws)
	defer Leave(uname)
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(chatroom.EVENT_MESSAGE, uname, string(p))
	}
}

func broadcastWebSocket(event chatroom.Event, subscribers list.list) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Println("Fail to marshal event:", err)
		return
	}
	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}
