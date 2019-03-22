package db

import (
	"fmt"
	"gopush/conf"
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
	dbType := mainConf.Postgresql.DbType
	dbUserName := mainConf.Postgresql.DbUserName
	dbUserPassword := mainConf.Postgresql.DbPassword
	dbName := mainConf.Postgresql.DbName
	dbPort := strconv.Itoa(int(mainConf.Postgresql.DbPort))
	dbHost := mainConf.Postgresql.DbHost
	//user=gorm password=gorm DB.name=gorm port=9920 sslmode=disable
	var err error
	var db *gorm.DB
	db, err = gorm.Open(dbType, "host="+ dbHost +" user=" + dbUserName + " password=" + dbUserPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable")
	fmt.Println("host="+ dbHost +" user=" + dbUserName + " password=" + dbUserPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable")
	if err != nil {
		panic(err)
	}
	factory := new(Session)
	factory.DB = db
	return factory
}
