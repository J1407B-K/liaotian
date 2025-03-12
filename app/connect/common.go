package connect

import (
	"github.com/gorilla/websocket"
	"sync"
)

// Clients 用于存储所有 WebSocket 客户端连接
var Clients = make(map[*websocket.Conn]string)
var ClientsMutex sync.Mutex
