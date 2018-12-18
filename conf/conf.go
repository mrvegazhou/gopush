package conf

type DbServer struct {
	DbType     string
	DbUserName string
	DbPassword string
	DbName     string
	Port       uint `default:"3306"`
}

type MainConfig struct {
	APPName   string `default:"app name"`
	Port      int    `default:8080`
	Addr      string `default:127.0.0.1`
	DbServers struct {
		DbServer1 DbServer
	}
	Jwt struct {
		Key string `default:"secret"`
	}
}
