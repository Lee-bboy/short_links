package router

import (
	"net/http"

	"short_links/common"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = common.Logger(route.HandlerFunc, route.Name)

		//API认证
		if route.Name != "jump to long link" {
			handler = common.Auth(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
