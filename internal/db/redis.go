package db

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	Client *redis.Client
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisClient{Client: rdb}
}

func (r *RedisClient) SetOTP(ctx context.Context, email string, otp string, duration time.Duration) error {
	key := "otp:" + email
	return r.Client.Set(ctx, key, otp, duration).Err()
}

func (r *RedisClient) GetOTP(ctx context.Context, email string) (string, error) {
	key := "otp:" + email
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisClient) DeleteOTP(ctx context.Context, email string) error {
	key := "otp:" + email
	return r.Client.Del(ctx, key).Err()
}