package sns

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewClient(acf aws.Config, baseEndpoint string) *sns.Client {
	client := sns.NewFromConfig(acf, func(o *sns.Options) {
		if baseEndpoint != "" {
			o.BaseEndpoint = aws.String(baseEndpoint)
			fmt.Println("Configurando SNS con LocalStack")
		} else {
			fmt.Println("Configurando SNS con AWS")
		}
	})

	return client
}
