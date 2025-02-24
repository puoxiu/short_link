package bloom

import (
	"context"
	"fmt"
	"hash/fnv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

// BloomFilter 布隆过滤器结构体
type BloomFilter struct {
    client       *redis.Client
    rs           *redsync.Redsync
    bloomFilterKey     string
    lockKey      string
    lockTimeout  time.Duration
}

// NewBloomFilter 创建一个新的布隆过滤器实例
func NewBloomFilter(RedisAddr, BloomFilterKey, LockKey string, LockTimeout time.Duration) (*BloomFilter) {
    client := redis.NewClient(&redis.Options{
        Addr: RedisAddr,
    })

    pool := goredis.NewPool(client)
    rs := redsync.New(pool)

    return &BloomFilter{
        client:       client,
        rs:           rs,
        bloomFilterKey:     BloomFilterKey,
        lockKey:      LockKey,
        lockTimeout:  LockTimeout,
    }
}

// Add 向布隆过滤器中添加一个元素
func (bf *BloomFilter) Add(ctx context.Context, item string) error {
	fmt.Println("添加Add item: ", item)
    mutex := bf.rs.NewMutex(bf.lockKey, redsync.WithExpiry(bf.lockTimeout))
    if err := mutex.LockContext(ctx); err != nil {
        return err
    }
    defer mutex.UnlockContext(ctx)

    hash1, hash2 := bf.hash(item)

    if err := bf.client.SetBit(ctx, bf.bloomFilterKey, hash1, 1).Err(); err != nil {
        return err
    }
    if err := bf.client.SetBit(ctx, bf.bloomFilterKey, hash2, 1).Err(); err != nil {
        return err
    }

    return nil
}

// Contains 检查布隆过滤器中是否包含某个元素
func (bf *BloomFilter) Contains(ctx context.Context, item string) (bool, error) {
    hash1, hash2 := bf.hash(item)

    bit1, err := bf.client.GetBit(ctx, bf.bloomFilterKey, hash1).Result()
    if err != nil {
        return false, err
    }
    bit2, err := bf.client.GetBit(ctx, bf.bloomFilterKey, hash2).Result()
    if err != nil {
        return false, err
    }
	if bit1 == 1 && bit2 == 1 {
		fmt.Println("布隆过滤器Contains item: ", item)
		return true, nil
	}

    return false, nil
}

// hash 计算元素的哈希值
// todo: 增加hash个数； 改进hash函数-->减小bloom过滤器的误判率
func (bf *BloomFilter) hash(item string) (int64, int64) {
    h1 := fnv.New64a()
    h1.Write([]byte(item))
    hash1 := h1.Sum64()

    h2 := fnv.New64a()
    h2.Write([]byte(item + "salt"))
    hash2 := h2.Sum64()

    // 限制位偏移量在合理范围内
    maxOffset := int64(1 << 32) - 1
    return int64(hash1 % uint64(maxOffset)), int64(hash2 % uint64(maxOffset))
}