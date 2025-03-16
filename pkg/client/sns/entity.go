package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Config struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
}

type Service interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}
