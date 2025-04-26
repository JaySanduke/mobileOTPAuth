package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewRedis(redisAddr string) *RedisClient {

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
		// Password: "yourpassword",
		DB:         0,
		PoolSize:   10,
		MaxRetries: 3,
	})
	ctx := context.Background()

	return &RedisClient{
		Client: rdb,
		Ctx:    ctx,
	}
}

func (r *RedisClient) SetOTP(mobile string, otp string, ttl time.Duration) error {
	return r.Client.Set(r.Ctx, "otp:"+mobile, otp, ttl).Err()
}

func (r *RedisClient) GetOTP(mobile string) (string, error) {
	return r.Client.Get(r.Ctx, "otp:"+mobile).Result()
}

func (r *RedisClient) DeleteOTP(mobile string) error {
	return r.Client.Del(r.Ctx, "otp:"+mobile).Err()
}
