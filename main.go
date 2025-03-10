package main

import (
	"fmt"
	"log"
	"websocket/config"
	"websocket/connect"
	"websocket/database"
	"websocket/router"
)

func main() {
	// 启动 Kafka 消费者协程
	go connect.ConsumeKafkaMessages()

	config.SetupViper()

	database.InitMysql()

	r := router.InitRouter()
	port := 8088
	log.Printf("WebSocket 聊天服务器启动")
	// 启动 HTTP 服务器
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
