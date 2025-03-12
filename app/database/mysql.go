package database

import "websocket/app/global"

func Create(value interface{}) error {
	return global.MysqlDB.Create(value).Error
}

func Select(query string, save interface{}, cond ...interface{}) error {
	return global.MysqlDB.Where(query, cond...).First(save).Error
}
