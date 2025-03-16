package sns

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Sns struct {
	Cliente Service
}

func NewClient(acf aws.Config, baseEndpoint string, l *logrus.Logger) *Sns {
	client := sns.NewFromConfig(acf, func(o *sns.Options) {
		if baseEndpoint != "" {
			o.BaseEndpoint = aws.String(baseEndpoint)
			l.Debug(fmt.Sprintf("Configurando SNS con LocalStack, endpoint %s", baseEndpoint))
		} else {
			l.Debug("Configurando SNS con AWS")
		}
	})

	return &Sns{Cliente: client}
}
