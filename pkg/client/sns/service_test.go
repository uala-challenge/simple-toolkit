package sns

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewSNSClient(t *testing.T) {
	cfg := Config{
		Endpoint: "http://localhost:4566",
	}
	awsCfg := aws.Config{}
	client := NewClient(awsCfg, cfg.Endpoint, logrus.New())
	assert.NotNil(t, client, "El cliente SNS no debería ser nil")
}

func TestNewSNSClientNoBaseEndpoint(t *testing.T) {
	cfg := Config{}
	awsCfg := aws.Config{}
	client := NewClient(awsCfg, cfg.Endpoint, logrus.New())
	assert.NotNil(t, client, "El cliente SNS no debería ser nil")
}
