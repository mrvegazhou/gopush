package db

import (
	"../../conf"
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strconv"
)

// Session 会话
type Session struct {
	DB	*gorm.DB
}

func ConnectDB(mainConf *conf.MainConfig) *Session {
	dbType := mainConf.DbType
	dbUserName := mainConf.DbUserName
	dbUserPassword := mainConf.DbPassword
	dbName := mainConf.DbName
	dbPort := strconv.Itoa(mainConf.Port)
	//user=gorm password=gorm DB.name=gorm port=9920 sslmode=disable
	var err error
	var db *gorm.DB
	db, err = gorm.Open(dbType, "user=" + dbUserName + " password=" + dbUserPassword + " DB.name=" + dbName + " port=" + dbPort + " sslmode=disable")
	if err != nil {
		panic(err)
	}
	factory := new(Session)
	factory.DB = db
	return factory
}
