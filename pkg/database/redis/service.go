package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *redis.Client
	config Config
	logger log.Service
}

var _ Service = (*service)(nil)

func NewService(cfg Config, logger log.Service) Service {
	options := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		DB:   cfg.DB,
	}

	if cfg.Password != "" {
		options.Password = cfg.Password
	}

	client := redis.NewClient(options)

	return &service{
		client: client,
		config: cfg,
		logger: logger,
	}
}

func (s *service) Set(ctx context.Context, key string, value any, expiration int) error {
	if s.client == nil {
		return errors.New(ErrorServiceNotEnabled)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	err := s.client.Set(ctx, key, value, time.Duration(expiration)*time.Hour).Err()
	if err != nil {
		s.logger.Error(ctx, err, "Error al establecer valor en Redis", map[string]interface{}{
			"key": key,
		})
		return err
	}

	s.logger.Info(ctx, "Valor establecido en Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}

func (s *service) Get(ctx context.Context, key string) (string, error) {
	if s.client == nil {
		return "", errors.New(ErrorServiceNotEnabled)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	value, err := s.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		s.logger.Warn(ctx, "Clave no encontrada en Redis", map[string]interface{}{
			"key": key,
		})
		return "", nil
	} else if err != nil {
		s.logger.Error(ctx, err, "Error al obtener valor de Redis", map[string]interface{}{
			"key": key,
		})
		return "", err
	}

	s.logger.Info(ctx, "Valor recuperado desde Redis", map[string]interface{}{
		"key": key,
	})
	return value, nil
}

func (s *service) Delete(ctx context.Context, key string) error {
	if s.client == nil {
		return errors.New(ErrorServiceNotEnabled)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.Timeout)*time.Second)
	defer cancel()

	err := s.client.Del(ctx, key).Err()
	if err != nil {
		s.logger.Error(ctx, err, "Error al eliminar clave en Redis", map[string]interface{}{
			"key": key,
		})
		return err
	}

	s.logger.Info(ctx, "Clave eliminada en Redis", map[string]interface{}{
		"key": key,
	})
	return nil
}
