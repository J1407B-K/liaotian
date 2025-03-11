package flag

import (
	"log"
	"websocket/app/global"
	"websocket/app/model"
)

func DatabaseAutoMigrate() {
	var err error

	err = global.MysqlDB.Set("gorm:table_option", "Engine=InnoDB").AutoMigrate(
		&model.User{},
	)

	if err != nil {
		log.Fatalf("AutoMigrate err:%v", err)
	}
	log.Println("AutoMigrate success")
}
