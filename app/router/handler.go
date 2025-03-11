package router

import (
	"github.com/gin-gonic/gin"
	"websocket/app/service"
)

func NewWsHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		service.WsHandler(c) // 调用 service 的 WsHandler
	}
}

func NewUserRegister() func(*gin.Context) {
	return func(c *gin.Context) {
		service.UserRegister(c)
	}
}

func NewUserLogin() func(*gin.Context) {
	return func(c *gin.Context) {
		service.UserLogin(c)
	}
}
