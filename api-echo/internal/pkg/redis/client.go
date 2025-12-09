package redis

import (
	"strconv"

	"github.com/redis/go-redis/v9"

	"filmogophery/internal/pkg/config"
)

func NewClient(conf *config.Config) *redis.Client {
	db, _ := strconv.Atoi(conf.Redis.DB)

	client := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Host + ":" + conf.Redis.Port,
		Password: conf.Redis.Password,
		DB:       db,
	})

	return client
}
