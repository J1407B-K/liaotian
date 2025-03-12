package database

import (
	"context"
	"websocket/app/global"
)

func InsertMongo(ctx context.Context, msg interface{}) error {
	_, err := global.MongoMsgCollection.InsertOne(ctx, msg)
	return err
}
