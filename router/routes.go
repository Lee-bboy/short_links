package router

import (
	"net/http"

	. "short_links/controllers"
)

const (
	JSON  string = "application/json;charset=utf-8"
	HTML  string = "text/html"
	PLAIN string = "text/plain"
)

type Route struct {
	Name        string //日志中显示的路由名称
	Method      string //GET PUT POST DELETE ...
	Pattern     string //对用的访问路径
	HandlerFunc http.HandlerFunc
	ContentType string //返回的数据类型"application/json;charset=utf-8" 或者 "text/html" 等等 ...
}

type Routes []Route

var routes = Routes{
	Route{
		"convert long links to short links",
		"POST",
		"/api/short_links",
		ControllerApi.ConvertLinks,
		JSON,
	},
	Route{
		"jump to long link",
		"GET",
		"/{shortlink}",
		ControllerApi.JumpToLongLink,
		HTML,
	},
}
