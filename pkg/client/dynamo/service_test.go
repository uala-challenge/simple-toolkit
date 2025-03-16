package dynamo

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	log "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestNewService(t *testing.T) {
	cfg := Config{
		Endpoint: "http://localhost:4566",
	}
	awsCfg := aws.Config{}
	mockLogger := log.NewService(t)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return(nil)
	service := NewClient(awsCfg, cfg, mockLogger)
	assert.NotNil(t, service, "El servicio DynamoDB no debería ser nil")
}

func TestNewServiceNoEndpoint(t *testing.T) {
	cfg := Config{}
	awsCfg := aws.Config{}
	mockLogger := log.NewService(t)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return(nil)
	service := NewClient(awsCfg, cfg, mockLogger)
	assert.NotNil(t, service, "El servicio DynamoDB no debería ser nil")
}
