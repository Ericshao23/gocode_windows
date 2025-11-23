package distributedlock

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisLock struct {
	client *redis.Client
}

func NewRedisLock(addr, password string, db int) *RedisLock {
	return &RedisLock{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisLock) AcquireLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	result, err := r.client.SetNX(ctx, lockInfo.key, lockInfo.value, lockInfo.expiration).Result()
	if err != nil {
		return false, err
	}
	return result, nil
}

func (r *RedisLock) ReleaseLock(ctx context.Context, lockInfo *DistributedLockInfo) (bool, error) {
	val, err := r.client.Get(ctx, lockInfo.key).Result()
	if err != nil {
		// 如果 key 不存在，说明锁已经被释放
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	if val == lockInfo.value {
		_, err = r.client.Del(ctx, lockInfo.key).Result()
		return err == nil, err
	}
	return false, nil
}

func (r *RedisLock) RenewLock(ctx context.Context, lockInfo *DistributedLockInfo) error {
	exists, err := r.client.Expire(ctx, lockInfo.key, lockInfo.expiration).Result()
	if err != nil {
		return err
	}
	if !exists {
		return ErrLockNotHeld
	}
	return nil
}

func (r *RedisLock) BuildServiceType() string {
	return "redis"
}
