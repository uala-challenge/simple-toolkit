package log

import (
	"context"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type service struct {
	Log *logrus.Logger
}

var _ Service = (*service)(nil)

func NewService(c Config, l *logrus.Logger) Service {
	l.Level = loggerLevel(c.Level)
	fileLog := c.Path
	file, err := os.OpenFile(fileLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		l.SetOutput(file)
	} else {
		l.Warn("No se pudo abrir el archivo de log, usando salida estándar")
	}
	l.Info("app started...")
	return &service{
		Log: l,
	}
}

func (l *service) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	l.Log.WithContext(ctx).WithFields(fields).Info(msg)
}

func (l *service) Error(ctx context.Context, err error, msg string, fields map[string]interface{}) {
	if err == nil {
		err = errors.New(msg)
	} else {
		err = l.WrapError(err, msg)
	}

	l.Log.WithContext(ctx).WithFields(fields).Error(err)
}

func (l *service) Debug(ctx context.Context, fields map[string]interface{}) {
	l.Log.WithContext(ctx).WithFields(fields).Debug()
}

func (l *service) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	l.Log.WithContext(ctx).WithFields(fields).Warn(msg)
}

func (l *service) FatalError(ctx context.Context, err error, fields map[string]interface{}) {
	l.Log.WithContext(ctx).WithFields(fields).Fatal(err)
}

func (l *service) WrapError(err error, msg string) error {
	if err == nil {
		return errors.New(msg)
	}
	return errors.Wrap(err, msg)
}

func loggerLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	case "debug":
		return logrus.DebugLevel
	case "trace":
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}
