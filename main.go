package main

import (
	"./framework/config"
	"./httpserver/routes"
	// "fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
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
	router.LoadHTMLGlob("./httpserver/views/*.html")
	router.StaticFS("/static", http.Dir("./httpserver/views/static"))
	routes.CreateRouter(router, Conf)
	http.ListenAndServe(":"+strconv.Itoa(Conf.Port), router)
}
