package rest

import (
	"context"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker/v2"
)

const (
	BackoffFactor                    = 2.0
	MaxJitterPercentage              = 0.5
	DefaultRetryCount         uint32 = 3
	DefaultRetryWaitTime             = 100 * time.Millisecond
	DefaultRetryMaxWaitTime          = 500 * time.Millisecond
	DefaultCBMaxRequests      uint32 = 5
	DefaultCBInterval                = 10 * time.Second
	DefaultCBTimeout                 = 5 * time.Second
	DefaultCBRequestThreshold uint32 = 5
	DefaultCBFailureRateLimit        = 50.0
)

type Config struct {
	TimeOut            time.Duration
	EnableLogging      bool
	BaseURL            string
	WithRetry          bool
	RetryCount         uint32
	RetryWaitTime      time.Duration
	RetryMaxWaitTime   time.Duration
	WithCB             bool
	CBName             string
	CBMaxRequests      uint32
	CBInterval         time.Duration
	CBTimeout          time.Duration
	CBRequestThreshold uint32
	CBFailureRateLimit float64
}

type Service interface {
	Get(ctx context.Context, endpoint string) (*resty.Response, error)
	Post(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error)
	Put(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error)
	Patch(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error)
	Delete(ctx context.Context, endpoint string) (*resty.Response, error)
	WithLogging(enable bool)
}

type client struct {
	baseURL   string
	requester *requester
	logger    *logrus.Logger
	logging   bool
}

type requester struct {
	httpClient *resty.Client
	breaker    *gobreaker.CircuitBreaker[any]
}
