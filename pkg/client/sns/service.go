package sns

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewClient(acf aws.Config, baseEndpoint string, l log.Service) *sns.Client {
	client := sns.NewFromConfig(acf, func(o *sns.Options) {
		if baseEndpoint != "" {
			o.BaseEndpoint = aws.String(baseEndpoint)
			l.Debug(context.Background(),
				map[string]interface{}{"message": "Configurando SNS con LocalStack", "endpoint": baseEndpoint})
		} else {
			l.Debug(context.Background(),
				map[string]interface{}{"message": "Configurando SNS con AWS"})
		}
	})

	return client
}
