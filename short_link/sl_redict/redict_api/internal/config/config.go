package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	ShortUrlDB struct {
		DSN string
	}

	BloomFilterKey	string
	BloomFilterRedisLockKey string
	BloomFilterLockTime int

	CacheRedis cache.CacheConf	// 添加缓存配置

}