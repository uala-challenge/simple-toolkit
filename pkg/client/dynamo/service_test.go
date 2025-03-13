package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	cfg := Config{
		TableName: "UalaChallenge",
		Endpoint:  "http://localhost:4566",
	}

	awsCfg := aws.Config{}
	service := NewClient(awsCfg, cfg)
	assert.NotNil(t, service, "El servicio DynamoDB no debería ser nil")
}
