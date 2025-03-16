package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	log "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		Host:    "127.9.9.9",
		Port:    6379,
		DB:      0,
		Timeout: 1,
	}

	mockLogger := log.NewService(t)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return(nil)
	cliente, err := NewClient(cfg, mockLogger)

	assert.NoError(t, err)
	assert.NotNil(t, cliente)
}

func TestNewClientWithPass(t *testing.T) {
	cfg := Config{
		Host:     "127.9.9.9",
		Port:     6379,
		DB:       0,
		Timeout:  1,
		Password: "pass",
	}
	mockLogger := log.NewService(t)
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return(nil)
	cliente, err := NewClient(cfg, mockLogger)

	assert.NoError(t, err)
	assert.NotNil(t, cliente)
}
