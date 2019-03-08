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
		DbDebug		bool  `default:false yaml:"dbDebug"`
	}
	Jwt struct {
		Key string `default:"secret" yaml:"key"`
	}
	Json struct{
		Pretty bool `default:true yaml:"pretty"`
	}
	Tcp struct{
		Address	string 	`default:"127.0.0.1" yaml:"address"`
		MaxConnCount int	`default:10 yaml:"maxConnCount"`
		AcceptCount int `default:10 yaml:"acceptCount"`
	}
}
