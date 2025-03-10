package model

type Config struct {
	RouterConfig
	MysqlConfig
}

type RouterConfig struct {
	Port string
}

type MysqlConfig struct {
	Username string
	Password string
	Addr     string
	DBName   string
}
