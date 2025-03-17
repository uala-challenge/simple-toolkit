package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Config struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type Service interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}
