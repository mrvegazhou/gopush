package db

import (
	"../../conf"
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectDB(mainConf *conf.MainConfig) (bool, error) {
	dbType := mainConf.DbServer.DbType
	dbUserName := mainConf.DbServer.DbUserName
	dbUserPassword := mainConf.DbServer.DbPassword
	dbName := mainConf.DbServer.DbName
	Port := mainConf.DbServer.Port

	DB, err = gorm.Open(dbType, "user="+dbUserName+" dbname="+dbName+" password="+dbUserPassword+" sslmode=disable")
	if err != nil {
		return false, err
	}
	return true, nil
}
