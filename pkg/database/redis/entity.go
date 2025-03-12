package redis

import (
	"context"
)

const (
	ErrorServiceNotEnabled = "redis Service not enabled"
)

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
	Timeout  int    `json:"timeout"`
}

type Service interface {
	Set(ctx context.Context, key string, value any, expiration int) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}
