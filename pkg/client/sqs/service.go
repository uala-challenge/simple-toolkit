package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func NewClient(acf aws.Config, cfg Config) *sqs.Client {
	return sqs.NewFromConfig(acf, func(o *sqs.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			fmt.Println("Configurando SQS con LocalStack")
		} else {
			fmt.Println("onfigurando SQS con AWS")
		}
	})
}
