package routes

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Method     string
	Pattern    string
	Handler    http.HandlerFunc
	Middleware mux.MiddlewareFunc
}

var routes []Route

func init() {
	//register("GET", "/movies", controllers.AllMovies, auth.TokenMiddleware)
}

func CreateRouter(conf interface{}) *mux.Router {
	router := mux.NewRouter()
	for _, route := range routes {
		r := router.Methods(route.Method).Path(route.Pattern)
		if route.Middleware != nil {
			r.Handler(route.Middleware(route.Handler, conf))
		} else {
			r.Handler(route.Handler)
		}
	}
	return router
}

func register(method, pattern string, handler http.HandlerFunc, middleware mux.MiddlewareFunc) {
	routes = append(routes, Route{method, pattern, handler, middleware})
}
