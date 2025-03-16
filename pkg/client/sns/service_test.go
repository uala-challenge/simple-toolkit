package sns

import (
	"testing"

	"github.com/stretchr/testify/mock"
	log "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewSNSClient(t *testing.T) {
	cfg := Config{
		Endpoint: "http://localhost:4566",
	}
	awsCfg := aws.Config{}
	mockLogger := log.NewService(t)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return(nil)
	client := NewClient(awsCfg, cfg.Endpoint, mockLogger)
	assert.NotNil(t, client, "El cliente SNS no debería ser nil")
}

func TestNewSNSClientNoBaseEndpoint(t *testing.T) {
	cfg := Config{}
	awsCfg := aws.Config{}
	mockLogger := log.NewService(t)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return(nil)
	client := NewClient(awsCfg, cfg.Endpoint, mockLogger)
	assert.NotNil(t, client, "El cliente SNS no debería ser nil")
}
