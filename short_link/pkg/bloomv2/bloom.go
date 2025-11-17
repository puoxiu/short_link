package bloomv2

import (
	"context"
	"fmt"
	"hash/fnv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/zeromicro/go-zero/core/logx"
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

func (bf *BloomFilter) Add(ctx context.Context, item string) error {
    // 获取分布式锁
    mutex := bf.rs.NewMutex(bf.lockKey, redsync.WithExpiry(bf.lockTimeout))
    if err := mutex.LockContext(ctx); err != nil {
        return err
    }

    // 看门狗机制
    renewTicker := time.NewTicker(bf.lockTimeout / 2)
    defer renewTicker.Stop()

    ctx2, cancel := context.WithCancel(ctx)
    defer cancel()

    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        for {
            select {
            case <-renewTicker.C:
                ok, err := mutex.ExtendContext(ctx2)
                if err != nil {
                    logx.Errorw("Failed to extend lock", logx.LogField{Key: "err", Value: err})
                    cancel() // 通知主协程锁续期失败
                    return
                }
                if!ok {
                    logx.Error("Lock extension failed, lock may have been released.")
                    cancel() // 通知主协程锁续期失败
                    return
                }
            case <-ctx2.Done():
                return
            }
        }
    }()

    // 确保 Goroutine 在主函数返回前退出
    defer wg.Wait()

    // 在 defer 中释放锁
    defer func() {
        // 使用两个变量接收返回值
        ok, err := mutex.UnlockContext(ctx2)
        if err != nil {
            logx.Errorw("Failed to unlock", logx.LogField{Key: "err", Value: err})
        }
        if!ok {
            logx.Error("Unlock failed, lock may have been released by other process.")
        }
    }()

    // 计算哈希值
    hash1, hash2 := bf.hash(item)

    // 使用 Redis 事务确保原子性
    pipe := bf.client.Pipeline()
    pipe.SetBit(ctx2, bf.bloomFilterKey, hash1, 1)
    pipe.SetBit(ctx2, bf.bloomFilterKey, hash2, 1)

    // 监听 ctx2.Done()，如果锁续期失败则提前返回
    select {
    case <-ctx2.Done():
        return fmt.Errorf("lock renewal failed")
    default:
        // 继续执行 Redis 操作
        if _, err := pipe.Exec(ctx2); err != nil {
            return err
        }
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
		// fmt.Println("布隆过滤器Contains item: ", item)
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


// Add 向布隆过滤器中添加一个元素
// func (bf *BloomFilter) Add(ctx context.Context, item string) error {
// 	// fmt.Println("添加Add item: ", item)
//     mutex := bf.rs.NewMutex(bf.lockKey, redsync.WithExpiry(bf.lockTimeout))
//     if err := mutex.LockContext(ctx); err != nil {
//         return err
//     }
//     defer mutex.UnlockContext(ctx)

// 	// 看门狗机制
// 	renewTicker := time.NewTicker(bf.lockTimeout / 2)
// 	defer renewTicker.Stop()

// 	ctx2, cancel := context.WithCancel(ctx)
// 	defer cancel()
//     go func() {
//         for {
//             select {
//             case <-renewTicker.C:
//                 ok, err := mutex.ExtendContext(ctx)
//                 if err != nil {
//                     // 处理续期失败的情况
// 					logx.Errorw("Failed to extend lock", logx.LogField{Key: "err", Value: err})
//                     return
//                 }
//                 if!ok {
//                     // 锁续期未成功，可能锁已被释放
// 					logx.Error("Lock extension failed, lock may have been released.")
//                     return
//                 }
//             case <-ctx2.Done():
//                 return
//             }
//         }
//     }()

//     hash1, hash2 := bf.hash(item)
//     if err := bf.client.SetBit(ctx, bf.bloomFilterKey, hash1, 1).Err(); err != nil {
//         return err
//     }
//     if err := bf.client.SetBit(ctx, bf.bloomFilterKey, hash2, 1).Err(); err != nil {
//         return err
//     }

//     return nil
// }