package dynamo

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Dynamo struct {
	Cliente Service
}

func NewClient(acf aws.Config, cfg Config, l *logrus.Logger) *Dynamo {
	client := dynamodb.NewFromConfig(acf, func(o *dynamodb.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			l.Debug(fmt.Sprintf("Configurando Dynamo con LocalStack, endpoint %s", cfg.Endpoint))
		} else {
			l.Debug("Configurando Dynamo con AWS")
		}
	})
	return &Dynamo{Cliente: client}
}
