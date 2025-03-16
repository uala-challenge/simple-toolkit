package redis

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg Config, l *logrus.Logger) (*redis.Client, error) {
	options := &redis.Options{
		Addr:        fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:          cfg.DB,
		DialTimeout: time.Duration(cfg.Timeout) * time.Second,
	}
	if cfg.Password != "" {
		options.Password = cfg.Password
	}
	client := redis.NewClient(options)

	l.Debug(fmt.Sprintf("Configurando Redis con LocalStack endpoint: %s:%d", cfg.Host, cfg.Port))
	return client, nil
}
