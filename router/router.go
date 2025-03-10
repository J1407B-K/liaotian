package router

import (
	"github.com/gin-gonic/gin"
	"websocket/utils"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/")
	{
		v1.POST("/register", NewUserRegister())
		v1.POST("/login", NewUserLogin())
	}

	v2 := r.Group("/chat")
	v2.Use(utils.JWTAuthMiddleware())
	{
		v2.GET("/ws", NewWsHandler())
	}
	return r
}
