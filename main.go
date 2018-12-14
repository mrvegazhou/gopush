package main

import (
	"./framework/config"
	"fmt"
	"net/http"
	"routes/routes"
)

var Conf = struct {
	APPName string `default:"app name"`
	Port    int    `default:8080`
	Addr    string `default:127.0.0.1`
	Mongodb struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true"`
		Port     uint   `default:"3306"`
	}

	Jwt struct {
		Key string `default:"secret"`
	}
}{}

func main() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Printf("err:", err)
	// 	}
	// }()
	Config.Load(&Conf, false, "conf/app.yml")
	router := gin.Default()
	createRouter(router, Conf)
	http.ListenAndServe(Conf.Port+":"+Conf.Addr, router)
}
