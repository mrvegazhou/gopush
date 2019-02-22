package conf

type Postgresql struct {
	DbType     string `default:"postgresql"`
	DbUserName string
	DbPassword string
	DbName     string
	Port       uint `default:"5432"`
}

type MainConfig struct {
	APPName   string `default:"app name"`
	Port      int    `default:8080`
	Addr      string `default:127.0.0.1`
	Postgresql
	Jwt struct {
		Key string `default:"secret"`
	}
}
