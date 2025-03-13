package redis

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg Config) (*redis.Client, error) {
	options := &redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:          cfg.DB,
		DialTimeout: time.Duration(cfg.Timeout) * time.Second,
	}
	if cfg.Password != "" {
		options.Password = cfg.Password
	}
	client := redis.NewClient(options)

	fmt.Println("conexión a Redis establecida correctamente")
	return client, nil
}
