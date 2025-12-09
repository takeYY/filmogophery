package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	IRedisService interface {
		Get(ctx context.Context, key string, dest interface{}) error
		Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
		Clear(ctx context.Context) error
	}

	redisService struct {
		client *redis.Client
	}
)

func NewRedisService(client *redis.Client) IRedisService {
	return &redisService{client: client}
}

func (s *redisService) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

func (s *redisService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return s.client.Set(ctx, key, data, expiration).Err()
}

func (s *redisService) Clear(ctx context.Context) error {
	return s.client.FlushDB(ctx).Err()
}
