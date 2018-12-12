package longpolling

import (
	"chatroom/chatroom"
	"encoding/json"
)

var (
	ErrTimeout        = errors.New("time out")
	ErrClientClosed   = errors.New("client request closed")
	ErrExpired        = errors.New("client expired")
	ErrTooManyClients = errors.New("too many clients")
)

func JoinRoom(user string) revel.Result {
	chatroom.Join(user)
	msg, err := json.MarshalIndent(msg, "", "")
	if nil != err {
		return err
	}
	return msg
}
