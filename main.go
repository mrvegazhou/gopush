package main

import (
	"gopush/conf"
	"gopush/httpserver/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopush/framework/config"
	"net/http"
	"strconv"
	"sync"
	"gopush/framework/db"
	"gopush/framework/db/imctx"
	"gopush/framework/im/connect"
)

var (
	Conf        = new(conf.MainConfig)
	once        sync.Once
)

func httpServer(ctx *imctx.Context) {
	// 创建一个不包含中间件的路由器
	router := gin.New()
	router.LoadHTMLGlob("./httpserver/views/*.html")
	router.StaticFS("/static", http.Dir("./httpserver/views/static"))
	routes.InitHandler(ctx, router)
	routes.CreateRouter(router, ctx)
	router.Run(":"+strconv.Itoa(Conf.Port))
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

	once.Do(func() {
		errs := Config.Load(&Conf, false, "conf/app.yml")
		fmt.Println(errs)
	})

	ctx := imctx.NewContext(db.ConnectDB(Conf), Conf)

	go httpServer(ctx)

	server := connect.NewTCPServer(ctx)
	server.Start(ctx)

}
