package log

import (
	"context"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestErrorLog_WithMessage(t *testing.T) {
	service := NewService(Config{Level: "error", Path: ""}, logrus.New())
	assert.NotPanics(t, func() {
		service.Error(context.Background(), errors.New("Error de prueba"), "Error de prueba", map[string]interface{}{"key": "value"})
	})
}

func TestErrorLog_WithoutError(t *testing.T) {
	service := NewService(Config{Level: "error", Path: ""}, logrus.New())
	assert.NotPanics(t, func() {
		service.Error(context.Background(), nil, "Mensaje de error", map[string]interface{}{"key": "value"})
	})

}

func TestDebugLog(t *testing.T) {
	service := NewService(Config{Level: "debug", Path: ""}, logrus.New())
	assert.NotPanics(t, func() {
		service.Debug(context.Background(), map[string]interface{}{"key": "value"})
	})

}

func TestInfoLog(t *testing.T) {
	service := NewService(Config{Level: "", Path: ""}, logrus.New())
	assert.NotPanics(t, func() {
		service.Info(context.Background(), "Mensaje de prueba", map[string]interface{}{"key": "value"})
	})
}

func TestWarnLog(t *testing.T) {
	service := NewService(Config{Level: "panic", Path: ""}, logrus.New())
	assert.NotPanics(t, func() {
		service.Warn(context.Background(), "Mensaje de prueba", map[string]interface{}{"key": "value"})
	})
}

func TestWrapError(t *testing.T) {
	service := NewService(Config{Level: "fatal", Path: ""}, logrus.New())
	err := service.WrapError(errors.New("Error original"), "Mensaje adicional")
	assert.EqualError(t, err, "Mensaje adicional: Error original")
}

func TestWrapErrorNil(t *testing.T) {
	service := NewService(Config{Level: "trace", Path: ""}, logrus.New())
	err := service.WrapError(nil, "Mensaje adicional")
	assert.EqualError(t, err, "Mensaje adicional")
}
