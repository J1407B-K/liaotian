package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	connect2 "websocket/app/connect"
	"websocket/app/global"
	"websocket/app/model"
	"websocket/app/utils"
)

func UserRegister(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := global.MysqlDB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func UserLogin(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := global.MysqlDB.Where("username = ?", user.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "generate token failed" + err.Error(),
		})
		log.Fatalf("generate token error: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":  0,
		"msg":   "login success",
		"token": token,
	})
}

// WsHandler 处理 WebSocket 连接
func WsHandler(c *gin.Context) {
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := global.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("升级 WebSocket 连接失败:", err)
		return
	}
	defer conn.Close()

	// 将新客户端加入连接列表
	connect2.ClientsMutex.Lock()
	connect2.Clients[conn] = true
	connect2.ClientsMutex.Unlock()
	log.Println("新 WebSocket 客户端已连接")

	// 循环读取客户端消息
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取 WebSocket 消息出错:", err)
			break
		}
		// 只处理文本消息
		if messageType != websocket.TextMessage {
			log.Println("非文本消息，忽略")
			continue
		}
		log.Printf("收到 WebSocket 消息: %s\n", string(message))
		// 将客户端消息发布到 Kafka
		if err = connect2.ProduceMessage(message); err != nil {
			log.Println("发送消息到 Kafka 出错:", err)
		}
	}

	// 客户端断开时，从连接列表中移除
	connect2.ClientsMutex.Lock()
	delete(connect2.Clients, conn)
	connect2.ClientsMutex.Unlock()
	log.Println("WebSocket 客户端已断开")
}
