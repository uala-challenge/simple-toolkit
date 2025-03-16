package dynamo

import (
	"context"
	"fmt"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewClient(acf aws.Config, cfg Config, l log.Service) *dynamodb.Client {
	client := dynamodb.NewFromConfig(acf, func(o *dynamodb.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			l.Debug(context.Background(),
				map[string]interface{}{"message": "Configurando Dynamo con LocalStack", "endpoint": cfg.Endpoint})
		} else {
			l.Debug(context.Background(),
				map[string]interface{}{"message": "Configurando Dynamo con AWS"})
		}
	})
	fmt.Println("cliente DynamoDB inicializado correctamente")
	return client
}
