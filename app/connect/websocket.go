package connect

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"log"
	"regexp"
	"strings"
	"time"
	"websocket/app/database"
	"websocket/app/global"
	"websocket/app/model"
)

// 敏感词列表
var sensitiveWords = []string{"毒品", "赌博", "违法"}

// broadcastMessage 将消息广播给所有已连接的 WebSocket 客户端(public)
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
func ProduceMessage(username, target string, timestamp time.Time, message []byte) error {
	// 创建 Kafka 生产者（writer），写完后关闭
	writer := initKafkaWriter()
	defer writer.Close() // 使用完关闭

	// 过滤敏感词 & 清理非法字符
	cleanedContent := filterSensitiveWords(cleanMessage(string(message)))

	// 组织消息数据
	msg := model.Message{
		Username:   username,
		TargetUser: target,
		Content:    cleanedContent,
		Timestamp:  timestamp,
	}

	// JSON 编码
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 将消息写入 Kafka
	return writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: messageBytes, // 消息内容
		},
	)
}

// consumeKafkaMessages 不断从 Kafka 中读取消息，并广播给所有 WebSocket 客户端
func ConsumeKafkaMessages() {
	reader := initKafkaReader()
	defer reader.Close()

	for {
		// 使用 context.Background()
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("读取 Kafka 消息出错:", err)
			continue // 继续读取下一条消息
		}

		// 解析 JSON 消息
		var message model.Message
		err = json.Unmarshal(msg.Value, &message)
		if err != nil {
			log.Println("解析 Kafka 消息失败:", err)
			continue
		}

		err = database.InsertMongo(context.TODO(), message)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("存入MongoDB成功")

		log.Printf("收到 Kafka 消息: %s\n", string(msg.Value))
		// 如果是私聊，发送给目标用户
		if message.TargetUser != "" {
			sendMessageToUser(message.TargetUser, msg.Value)
		} else {
			// 否则广播消息
			broadcastMessage(msg.Value)
		}
	}
}

// kafka辅助函数
func initKafkaWriter() *kafka.Writer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(global.KafkaBroker), // 设置 Kafka broker 地址
		Topic:    global.KafkaTopic,             // 设置 Kafka topic
		Balancer: &kafka.LeastBytes{},           // 设置负载均衡策略
	}
	return writer
}

func initKafkaReader() *kafka.Reader {
	// 创建 Kafka 消费者（reader）
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{global.KafkaBroker},
		Topic:   global.KafkaTopic,
		GroupID: "chat-group",
	})
	return reader
}

// 过滤相关
func filterSensitiveWords(content string) string {
	for _, word := range sensitiveWords {
		re := regexp.MustCompile(regexp.QuoteMeta(word))
		content = re.ReplaceAllString(content, strings.Repeat("*", len(word)))
	}
	return content
}

func cleanMessage(content string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9一-龥\s]`)
	return re.ReplaceAllString(content, "")
}

// 私聊
// sendMessageToUser 发送私聊消息给指定用户
func sendMessageToUser(targetUser string, message []byte) {
	ClientsMutex.Lock()
	defer ClientsMutex.Unlock()

	// 找到目标用户对应的 WebSocket 连接
	for client, username := range Clients {
		if username == targetUser {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("向目标客户端写入消息失败，移除该连接:", err)
				client.Close()
				delete(Clients, client)
			}
			return
		}
	}
	log.Println("目标用户未在线:", targetUser)
}
