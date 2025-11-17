package localcache

import (
	"github.com/allegro/bigcache/v3"
)

type LocalCacheHundler struct {
	cache *bigcache.BigCache
}

func NewLocalCacheHundler(cache *bigcache.BigCache) *LocalCacheHundler {
	return &LocalCacheHundler{cache: cache}
}

// Set 写入缓存（key: 短链接，value: 长链接）
func (h *LocalCacheHundler) Set(shortCode, longURL string) error {	
	return h.cache.Set(shortCode, []byte(longURL))
}

// Get 读取缓存（返回长链接，未命中返回空字符串）
func (h *LocalCacheHundler) Get(shortCode string) (string, error) {
	val, err := h.cache.Get(shortCode)
	if err == bigcache.ErrEntryNotFound {
		return "", nil
	}
	return string(val), err
}

// Delete 删除缓存
func (h *LocalCacheHundler) Delete(shortCode string) error {
	return h.cache.Delete(shortCode)
}