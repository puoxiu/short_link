package config

import "github.com/zeromicro/go-zero/rest"


type Config struct {
	rest.RestConf

	UserDB struct {
		DSN string
	}
	
	Etcd string
	Auth Auth
	WhiteList []string	// 白名单
}

type Auth struct {
	AccessSecret string
	AccessExpire int
}
