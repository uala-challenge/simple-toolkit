package sqs

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

func TestNewSNSClient(t *testing.T) {
	cfg := Config{
		BaseEndpoint: "http://localhost:4566",
	}
	awsCfg := aws.Config{}
	client := NewClient(awsCfg, cfg)
	assert.NotNil(t, client, "El cliente SQS no debería ser nil")
}

func TestNewSNSClientNoBaseEndpoint(t *testing.T) {
	cfg := Config{}
	awsCfg := aws.Config{}
	client := NewClient(awsCfg, cfg)
	assert.NotNil(t, client, "El cliente SQS no debería ser nil")
}
