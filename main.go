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
	router := gin.Default()
	router.LoadHTMLGlob("./httpserver/views/*.html")
	router.StaticFS("/static", http.Dir("./httpserver/views/static"))
	routes.InitHandler(ctx, router)
	routes.CreateRouter(router, ctx)
	http.ListenAndServe(":"+strconv.Itoa(Conf.Port), router)
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
