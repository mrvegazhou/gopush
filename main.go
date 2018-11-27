package main

import (
	"./framework/config"
	"fmt"
	"net/http"
)

var Conf = struct {
	APPName string `default:"app name"`

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

	r := routes.CreateRouter(Conf)
	http.ListenAndServe(":8080", r)
}
