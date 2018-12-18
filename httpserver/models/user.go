package models

import (
	_ "encoding/json"
	"github.com/jinzhu/gorm"
)

type UserType int

const (
	Human UserType = iota
	Robot
)

type User struct {
	gorm.Model
	UserId       string   `json:"user_id"`
	UserName     string   `json:"user_name"`
	UserPassword string   `json:"user_password"`
	UserAvatar   string   `json:"user_avatar"`
	UserPhone    string   `json:"user_phone"`
	UserEmail    string   `json:"user_email"`
	UserGender   string   `json:"user_gender"`
	UserBirth    string   `json:"user_birth"`
	UserCity     string   `json:"user_city"`
	UserType     UserType `json:"user_type"`
}
