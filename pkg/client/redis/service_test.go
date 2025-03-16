package redis

import (
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	cfg := Config{
		Host:    "127.9.9.9",
		Port:    6379,
		DB:      0,
		Timeout: 1,
	}

	cliente, err := NewClient(cfg, logrus.New())

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
	cliente, err := NewClient(cfg, logrus.New())

	assert.NoError(t, err)
	assert.NotNil(t, cliente)
}
