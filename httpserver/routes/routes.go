package routes

import (
	"github.com/gin-gonic/gin"
	"../controller/chat"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route
//https://github.com/winlion/restgo/blob/master/controller/PageController.go
func init() {
	register("GET", "/chat", controllers., auth.TokenMiddleware)
	register("GET", "/movies", controllers., auth.TokenMiddleware)

}

func registerRouter(router *gin.Engine, conf interface{}) {
	new(controller.PageController).Router(router)
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}
