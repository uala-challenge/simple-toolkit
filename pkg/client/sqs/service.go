package sqs

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewClient(acf aws.Config, cfg Config, l log.Service) *sqs.Client {
	return sqs.NewFromConfig(acf, func(o *sqs.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			l.Debug(context.Background(),
				map[string]interface{}{"message": "Configurando SQS con LocalStack", "endpoint": cfg.Endpoint})
		} else {
			l.Debug(context.Background(),
				map[string]interface{}{"message": "Configurando SQS con AWS"})
		}
	})
}
