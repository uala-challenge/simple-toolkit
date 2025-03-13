package dynamo

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewClient(acf aws.Config, cfg Config) *dynamodb.Client {
	client := dynamodb.NewFromConfig(acf, func(o *dynamodb.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})
	fmt.Println("cliente DynamoDB inicializado correctamente")
	return client
}
