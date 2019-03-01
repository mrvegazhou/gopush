package conf

type MainConfig struct {
	APPName   string `default:"app name" yaml:"appname"`
	Port      int    `default:8080 yaml:"port"`
	Addr      string `default:"127.0.0.1" yaml:"addr"`
	Postgresql struct {
		DbType     string `default:"postgres" yaml:"dbType"`
		DbUserName string `yaml:"dbUserName"`
		DbPassword string `yaml:"dbPassword"`
		DbName     string `yaml:"dbName"`
		DbPort       uint `default:"5432" yaml:"dbPort"`
		DbHost 	   string `default:"127.0.0.1"  yaml:"dbHost"`
	}
	Jwt struct {
		Key string `default:"secret" yaml:"key"`
	}
	Json struct{
		Pretty bool `default:true yaml:"pretty"`
	}
}
