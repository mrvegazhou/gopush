package main

import (
	"./framework/config"
	"./framework/helper"
	"./framework/log"
	"./httpserver/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"sync"
)

type DbServer struct {
	DbType     string
	DbUserName string
	DbPassword string
	Port       uint `default:"3306"`
}

var (
	Conf = struct {
		APPName   string `default:"app name"`
		Port      int    `default:8080`
		Addr      string `default:127.0.0.1`
		DbServers struct {
			DbServer1 DbServer
		}
		Jwt struct {
			Key string `default:"secret"`
		}
	}{}
	once        sync.Once
	innerLogger *logger.InnerLogger
)

func init() {
	innerLogger = logger.GetInnerLogger()
}

func main() {
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		strLog := "longweb:main recover error => " + fmt.Sprintln(err)
	// 		os.Stdout.Write([]byte(strLog))
	// 		errMsg := strLog + string(debug.Stack())
	// 		innerLogger.Error(errMsg)
	// 		os.Stdout.Write([]byte(errMsg))
	// 	}
	// }()
	baseDir := helper.GetCurrentDirectory()
	once.Do(func() {
		Config.Load(&Conf, false, "conf/app.yml")
		logger.StartInnerLogHandler(baseDir)
	})
	fmt.Printf("Conf:", Conf)
	router := gin.Default()

	// 设置表单验证 validator.v9
	binding.Validator = new(utils.ValidatorV9)

	router.LoadHTMLGlob("./httpserver/views/*.html")
	router.StaticFS("/static", http.Dir("./httpserver/views/static"))
	routes.CreateRouter(router, Conf)
	http.ListenAndServe(":"+strconv.Itoa(Conf.Port), router)
}
