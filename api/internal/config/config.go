package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	Mysql Mysql
	Redis Redis
}

type Mysql struct {
	DataSource string
}
type Redis struct {
	Addr     string
	Password string
	DB       int
}