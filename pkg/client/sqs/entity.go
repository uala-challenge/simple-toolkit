package sqs

import (
	"context"
	"errors"
)

var ErrNoMessages = errors.New("no hay mensajes en la cola")

type Service interface {
	ReceiveMessage(ctx context.Context) (Message, error)
	DeleteMessage(ctx context.Context, receiptHandle string) error
}

type Message struct {
	ID            string
	Body          string
	ReceiptHandle string
}

type Config struct {
	QueueURL        string `json:"queue_url"`
	MaxRetries      int    `json:"max_retries"`
	MaxMessages     int32  `json:"max_messages"`
	WaitTimeSeconds int32  `json:"wait_time_seconds"`
	TimeoutSeconds  int    `json:"timeout_seconds"`
}
