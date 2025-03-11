package init

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"websocket/app/global"
)

func InitMysql() {
	dsn := global.Config.MysqlConfig.Username + ":" + global.Config.MysqlConfig.Password + "@tcp(" + global.Config.MysqlConfig.Addr + ")/" + global.Config.MysqlConfig.DBName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql error: %v", err)
	}
	global.MysqlDB = db
}
