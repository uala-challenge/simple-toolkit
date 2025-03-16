package sqs

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewClient(acf aws.Config, cfg Config, l *logrus.Logger) *sqs.Client {
	return sqs.NewFromConfig(acf, func(o *sqs.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			l.Debug(fmt.Sprintf("Configurando SQS con LocalStack, endpoint %s", cfg.Endpoint))
		} else {
			l.Debug("Configurando SQS con AWS")
		}
	})
}
