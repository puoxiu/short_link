package svc

import (
	"short_link_pro/pkg/bloom"
	"short_link_pro/shorturlmapmodel"
	"short_link_pro/sl_convert/convert_api/internal/config"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"short_link_pro/sl_convert/convert_api/sequence"
	"time"

)

type ServiceContext struct {
	Config config.Config

	ShortUrlMapModel shorturlmapmodel.ShortUrlMapModel
	Sequence sequence.Sequence
	ShortUrlBlackList map[string]struct{}

	// bloom filter
	BloomFilter *bloom.BloomFilter
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.ShortUrlDB.DSN)
	// 将黑名单数据加载到 map 中
	m := make(map[string]struct{}, len(c.ShortUrlBlackList))
	for _, v := range c.ShortUrlBlackList {
		m[v] = struct{}{}
	}
	return &ServiceContext{
		Config: c,
		ShortUrlMapModel: shorturlmapmodel.NewShortUrlMapModel(conn, c.CacheRedis),
		Sequence: sequence.NewMySQL(c.SequenceDB.DSN),
		ShortUrlBlackList: m,
		BloomFilter: bloom.NewBloomFilter(c.CacheRedis[0].Host, c.BloomFilterKey, c.BloomFilterRedisLockKey, time.Duration(c.BloomFilterLockTime)*time.Second),
	}
}