package init

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
	"websocket/app/global"
)

func ConnectMongoDB() {
	//设置连接
	clientOptions := options.Client().ApplyURI(global.Config.MongoConfig.Addr)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	//连接
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("mongo connect failed, err:%v\n", err)
	}

	//测试连接
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("mongo ping failed, err:%v\n", err)
	}

	global.Mongo = client
	log.Println("mongo connect success")
}
