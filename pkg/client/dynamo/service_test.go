package dynamo

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	cfg := Config{
		Endpoint: "http://localhost:4566",
	}
	awsCfg := aws.Config{}
	service := NewClient(awsCfg, cfg, logrus.New())
	assert.NotNil(t, service, "El servicio DynamoDB no debería ser nil")
}

func TestNewServiceNoEndpoint(t *testing.T) {
	cfg := Config{}
	awsCfg := aws.Config{}
	service := NewClient(awsCfg, cfg, logrus.New())
	assert.NotNil(t, service, "El servicio DynamoDB no debería ser nil")
}
