package model

type Config struct {
	RouterConfig
	MysqlConfig
	MongoConfig
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

type MongoConfig struct {
	Addr string
}
