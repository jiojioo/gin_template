// Package redis owns the global Redis client.
package redis

import (
	"context"

	"github.com/jiojioo/gin_template/internal/config"
	redislib "github.com/redis/go-redis/v9"
)

var Client *redislib.Client

func Init(cfg config.RedisConfig) error {
	client := redislib.NewClient(&redislib.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		return err
	}
	Client = client
	return nil
}
