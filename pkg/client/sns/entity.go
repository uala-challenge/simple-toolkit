package sns

import (
	"context"
)

type Service interface {
	Accept(ctx context.Context, message interface{}) error
}

type Config struct {
	TopicARN       string `json:"topic_arn"`
	MaxRetries     int    `json:"max_retries"`
	TimeoutSeconds int    `json:"timeout_seconds"`
}
