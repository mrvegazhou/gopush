package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Group struct {
	gorm.Model
	GroupName       string    `json:"group_name"`
	GroupCreateTime time.Time `json:"group_create_time"`
	GroupMembers    []User    `json:"group_members"`
	GroupAddr       string    `json:"group_addr"`
}
