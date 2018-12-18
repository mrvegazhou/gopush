package db

import (
	"../../conf"
	_ "database/sql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func ConnectDB(mainConf *conf.MainConfig) bool {

}
