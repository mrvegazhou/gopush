package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Friends struct {
	gorm.Model
	UserAId           int
	UserBId           int
	FriendsCreateTime time.Time `json:"friends_create_time"`
}
