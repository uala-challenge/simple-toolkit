package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg Config, l log.Service) (*redis.Client, error) {
	options := &redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:          cfg.DB,
		DialTimeout: time.Duration(cfg.Timeout) * time.Second,
	}
	if cfg.Password != "" {
		options.Password = cfg.Password
	}
	client := redis.NewClient(options)

	l.Debug(context.Background(),
		map[string]interface{}{"message": "Configurando Redis con LocalStack",
			"endpoint": fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)})
	return client, nil
}
