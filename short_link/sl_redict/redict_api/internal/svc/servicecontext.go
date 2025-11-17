package svc

import (
	"context"
	"short_link_pro/pkg/bloomv2"
	"short_link_pro/shorturlmapmodel"
	"short_link_pro/sl_redict/redict_api/internal/config"
	"short_link_pro/sl_redict/redict_api/internal/middleware"
	"short_link_pro/sl_redict/redict_models"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"short_link_pro/pkg/localcache"
)

type ServiceContext struct {
	Config             config.Config
	ClientIPMiddleware rest.Middleware

	ShortUrlMapModel shorturlmapmodel.ShortUrlMapModel
	ShortUrlAccessLogModel redict_models.ShortUrlAccessLogModel
	// bloom filter
	BloomFilter *bloomv2.BloomFilter
	// local cache
	LocalCacheHundler *localcache.LocalCacheHundler
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)

	bigcache, err := bigcache.New(
		context.Background(),
		bigcache.Config{
			Shards:             256,	// 分片数
			LifeWindow:         c.BigCache.Expiration,
			MaxEntriesInWindow: c.BigCache.MaxEntries,	// 
			HardMaxCacheSize:   1024,	// 最大缓存1024MB
			Verbose:            false,	// 关闭详细日志
		},
	)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:             c,
		ClientIPMiddleware: middleware.NewClientIPMiddleware().Handle,

		ShortUrlMapModel: shorturlmapmodel.NewShortUrlMapModel(conn, c.CacheRedis),
		ShortUrlAccessLogModel: redict_models.NewShortUrlAccessLogModel(conn),
		BloomFilter: bloomv2.NewBloomFilter(c.CacheRedis[0].Host, c.BloomFilterKey, c.BloomFilterRedisLockKey, time.Duration(c.BloomFilterLockTime)*time.Second),
		LocalCacheHundler: localcache.NewLocalCacheHundler(bigcache),
	}
}
