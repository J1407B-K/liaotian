package connect

import (
	"context"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"log"
	"websocket/app/global"
)

// broadcastMessage 将消息广播给所有已连接的 WebSocket 客户端
func broadcastMessage(message []byte) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()
	for client := range Clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("向客户端写入消息失败，移除该连接:", err)
			client.Close()
			delete(Clients, client)
		}
	}
}

// produceMessage 将收到的消息发送到 Kafka
func ProduceMessage(message []byte) error {
	// 创建 Kafka 生产者（writer），写完后关闭
	writer := &kafka.Writer{
		Addr:     kafka.TCP(global.KafkaBroker), // 设置 Kafka broker 地址
		Topic:    global.KafkaTopic,             // 设置 Kafka topic
		Balancer: &kafka.LeastBytes{},           // 设置负载均衡策略
	}
	defer writer.Close() // 使用完关闭

	// 将消息写入 Kafka
	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: message, // 消息内容
		},
	)
}

// consumeKafkaMessages 不断从 Kafka 中读取消息，并广播给所有 WebSocket 客户端
func ConsumeKafkaMessages() {
	// 创建 Kafka 消费者（reader）
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{global.KafkaBroker},
		Topic:   global.KafkaTopic,
		GroupID: "chat-group",
	})
	defer reader.Close()

	for {
		// 使用 context.Background() 而非 nil
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("读取 Kafka 消息出错:", err)
			continue // 继续读取下一条消息
		}
		log.Printf("收到 Kafka 消息: %s\n", string(msg.Value))
		broadcastMessage(msg.Value)
	}
}
