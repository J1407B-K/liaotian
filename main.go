package main

import (
	"fmt"
	"log"
	"websocket/app/config"
	"websocket/app/connect"
	"websocket/app/flag"
	"websocket/app/initialize"
	"websocket/app/router"
)

func main() {
	// 启动 Kafka 消费者协程
	go connect.ConsumeKafkaMessages()

	//初始化配置读取
	config.SetupViper()

	//初始化数据库
	initialize.InitMysql()
	initialize.ConnectMongoDB()

	//自动建表
	option := flag.Parse()
	if flag.IsWebStop(option) {
		flag.SwitchOption(option)
	}

	r := router.InitRouter()
	port := 8088
	log.Printf("WebSocket 聊天服务器启动")
	// 启动 HTTP 服务器
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}
