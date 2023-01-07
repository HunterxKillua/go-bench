package server

import (
	"ginBlog/api"
	"os"
	md "ginBlog/middlewares"
	"github.com/gin-gonic/gin"
)

type Api struct {
	Instance *gin.Engine
	Group    map[string]*gin.RouterGroup
}

/* 启动服务 */
func Run(port string) {
	instance := gin.Default()
	instance.Use(md.InterceptRequest())
	for _, group := range api.Groups {
		group(instance)
	}
	for _, route := range api.Routes {
		route(instance)
	}
	if err := instance.Run(port); err != nil {
		os.Exit(-1)
	}
}
