package model

type User struct {
	UserId   int64  `gorm:"primary_key;auto_increment" json:"userId"`
	Username string `gorm:"unique;not null;size:50" json:"username"`
	Password string `gorm:"size:255;not null"  json:"password"`
}
