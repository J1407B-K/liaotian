package global

import (
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"net/http"
	"websocket/app/model"
)

const (
	// KafkaBroker 地址
	KafkaBroker = "127.0.0.1:9092"

	// KafkaTopic Topic名称，请确保该主题已创建
	KafkaTopic = "local-dev"
)

// 初始设置
var (
	MysqlDB *gorm.DB
	Config  *model.Config
	Mongo   *mongo.Client
)

// Upgrader 创建一个 WebSocket 升级器，允许所有跨域请求
var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
