package models

import (
	_ "encoding/json"
	"github.com/jinzhu/gorm"
)

type SendType int

type MsgType int

const (
	// one to one
	O2O SendType = iota
	// one to many, for group chat
	O2M
	// many to one, for official account receive
	M2O
)

const (
	Text MsgType = iota
	Image
	Video
	UrlLink
	Contact
	System
	Notification
)

type Msg struct {
	gorm.Model
	SenderId   int      `json:"sender_id"`
	TargetId   int      `json:"target_id"`
	SendTime   string   `json:"send_time"`
	ReadTime   string   `json:"read_time"`
	MsgContent string   `json:"msg_content"`
	SendType   SendType `json:"send_type"`
	MsgType    int      `json:"msg_type"`
}
