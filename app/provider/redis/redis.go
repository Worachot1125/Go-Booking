package provider

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisProvider struct {
	Client *redis.Client
}

func NewRedisProvider(options *redis.Options) *RedisProvider {
	client := redis.NewClient(options)
	return &RedisProvider{Client: client}
}

func (r *RedisProvider) Ping(ctx context.Context) error {
	_, err := r.Client.Ping(ctx).Result()
	return err
}

func (r *RedisProvider) Close() error {
	return r.Client.Close()
}
