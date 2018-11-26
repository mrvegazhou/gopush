package main

import (
	"./framework/config"
	"fmt"
)

var Conf = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}
}{}

func main() {
	defer func() {
		if err := recover(); err != nil {
		}
	}()
	Config.Load(&Conf, false, "conf/app.yml")
	fmt.Printf("config: %#v", Conf)
}
