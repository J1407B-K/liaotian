package model

import "time"

type Message struct {
	Username   string    `json:"username" bson:"username"`       // 用户名
	Content    string    `json:"content" bson:"content"`         // 处理后的消息内容（经过敏感词过滤等处理）
	Timestamp  time.Time `json:"timestamp" bson:"timestamp"`     // 消息时间戳
	TargetUser string    `json:"target_user" bson:"target_user"` // 目标用户（私聊时使用）
}
