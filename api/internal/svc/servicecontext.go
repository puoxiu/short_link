package svc

import (
	"short_link_svc/api/internal/config"
	"short_link_svc/core"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config

	//
	Mysql  *gorm.DB
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysql := core.InitGorm(c.Mysql.DataSource)
	redis := core.InitRedis(c.Redis.Addr, c.Redis.Password, c.Redis.DB)
	return &ServiceContext{
		Config: c,
		Mysql: mysql,
		Redis: redis,
	}
}

