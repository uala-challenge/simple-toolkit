package rest

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker/v2"
)

var _ Service = (*client)(nil)

func NewClient(cfg Config, l log.Service) Service {
	setDefaultConfig(&cfg)
	r := &requester{
		httpClient: createHttpClient(cfg, l, cfg.TimeOut),
		breaker:    createCB(cfg, l),
	}

	return &client{
		baseURL:   cfg.BaseURL,
		requester: r,
		logger:    l,
		logging:   cfg.EnableLogging,
	}
}

func setDefaultConfig(cfg *Config) {
	defaults := map[*uint32]uint32{
		&cfg.RetryCount:         DefaultRetryCount,
		&cfg.CBMaxRequests:      DefaultCBMaxRequests,
		&cfg.CBRequestThreshold: DefaultCBRequestThreshold,
	}
	for k, v := range defaults {
		if *k == 0 || *k > 100 {
			*k = v
		}
	}

	defaultsFloat := map[*float64]float64{
		&cfg.CBFailureRateLimit: DefaultCBFailureRateLimit,
	}
	for k, v := range defaultsFloat {
		if *k <= 0 || *k > 100 {
			*k = v
		}
	}

	defaultsDuration := map[*time.Duration]time.Duration{
		&cfg.RetryWaitTime:    DefaultRetryWaitTime,
		&cfg.RetryMaxWaitTime: DefaultRetryMaxWaitTime,
		&cfg.CBInterval:       DefaultCBInterval,
		&cfg.CBTimeout:        DefaultCBTimeout,
	}
	for k, v := range defaultsDuration {
		if *k <= 0 || *k > time.Minute {
			*k = v
		}
	}
}

func checkBreakerState(counts gobreaker.Counts, c Config, l log.Service) bool {
	var failureRate float64
	if counts.Requests > 0 {
		failureRate = float64(counts.TotalFailures) / float64(counts.Requests)
	}
	l.Info(context.Background(), "Circuit Breaker Metrics",
		map[string]interface{}{
			"Total Requests":     counts.Requests,
			"Total Successes":    counts.TotalSuccesses,
			"Total Failures":     counts.TotalFailures,
			"Failure Rate":       failureRate,
			"ConsecutiveFails":   counts.ConsecutiveFailures,
			"ConsecutiveSuccess": counts.ConsecutiveSuccesses,
		})
	if counts.ConsecutiveFailures > c.CBMaxRequests || (counts.Requests >= c.CBRequestThreshold && failureRate > c.CBFailureRateLimit) {
		l.Info(context.Background(), "Circuit Breaker se abrirá debido a una alta tasa de fallos.", nil)
		return true
	}
	return false
}

func createCB(c Config, l log.Service) *gobreaker.CircuitBreaker[any] {
	if !c.WithCB {
		return nil
	}
	cbConfig := gobreaker.Settings{
		Name:        c.CBName,
		MaxRequests: c.CBMaxRequests,
		Interval:    c.CBInterval,
		Timeout:     c.CBTimeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return checkBreakerState(counts, c, l)
		},

		OnStateChange: func(name string, from, to gobreaker.State) {
			l.Warn(context.Background(), "Circuit Breaker state changed", map[string]interface{}{
				"name": name,
				"from": from,
				"to":   to,
			})
		},
	}

	cb := gobreaker.NewCircuitBreaker[any](cbConfig)
	return cb
}

func exponentialBackoffWithJitter(initialWaitTime, maxWaitTime time.Duration, attempt int, l log.Service) time.Duration {
	if attempt <= 0 {
		attempt = 1
	}

	baseWaitTime := initialWaitTime * time.Duration(math.Pow(BackoffFactor, float64(attempt-1)))
	jitter := time.Duration(rand.Float64() * float64(baseWaitTime) * MaxJitterPercentage)
	waitTime := baseWaitTime + jitter

	if waitTime > maxWaitTime {
		waitTime = maxWaitTime
	}

	l.Debug(context.Background(), map[string]interface{}{
		"attempt":   attempt,
		"baseTime":  baseWaitTime,
		"jitter":    jitter,
		"waitTime":  waitTime,
		"maxWait":   maxWaitTime,
		"factor":    BackoffFactor,
		"jitterPct": MaxJitterPercentage,
	})

	return waitTime
}

func retryAfterFunc(initialWaitTime, maxWaitTime time.Duration, l log.Service) func(*resty.Client, *resty.Response) (time.Duration, error) {
	return func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
		attempt := resp.Request.Attempt
		return exponentialBackoffWithJitter(initialWaitTime, maxWaitTime, attempt, l), nil
	}
}

func createHttpClient(c Config, l log.Service, timeout time.Duration) *resty.Client {
	client := resty.New()
	if timeout > 0 {
		client.SetTimeout(timeout)
	}
	if c.WithRetry {
		client.SetRetryCount(int(c.RetryCount)).
			SetRetryAfter(retryAfterFunc(c.RetryWaitTime, c.RetryMaxWaitTime, l)).
			AddRetryCondition(func(r *resty.Response, err error) bool {
				return err != nil || r.StatusCode() >= 500
			})
	}
	return client
}

func (c *client) WithLogging(enable bool) {
	c.logging = enable
}

func (c *client) executeRequest(ctx context.Context, reqFunc func(ctx context.Context) (*resty.Response, error)) (*resty.Response, error) {
	ctx, cancel := c.ensureContextWithTimeout(ctx, 10*time.Second)
	defer cancel()

	if c.requester.breaker != nil {
		return c.executeWithCircuitBreaker(ctx, reqFunc)
	}

	return c.executeWithoutCircuitBreaker(ctx, reqFunc)
}

func (c *client) ensureContextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if _, hasDeadline := ctx.Deadline(); hasDeadline {
		return ctx, func() {
			c.logger.Info(ctx, "timeout defined", nil)
		}
	}
	return context.WithTimeout(ctx, timeout)
}

func (c *client) executeWithCircuitBreaker(ctx context.Context, reqFunc func(ctx context.Context) (*resty.Response, error)) (*resty.Response, error) {
	result, err := c.requester.breaker.Execute(func() (interface{}, error) {
		return c.performRequest(ctx, reqFunc)
	})

	if err != nil {
		return c.handleCircuitBreakerError(ctx, err)
	}

	return result.(*resty.Response), nil
}

func (c *client) executeWithoutCircuitBreaker(ctx context.Context, reqFunc func(ctx context.Context) (*resty.Response, error)) (*resty.Response, error) {
	resp, err := reqFunc(ctx)
	if err != nil {
		return nil, err
	}
	return c.validateAndReturnResponse(resp)
}

func (c *client) performRequest(ctx context.Context, reqFunc func(ctx context.Context) (*resty.Response, error)) (*resty.Response, error) {
	resp, err := reqFunc(ctx)
	if err != nil {
		c.logRequestFailure(ctx, err)
		return nil, err
	}
	if err := validateResponse(resp); err != nil {
		c.logHttpError(ctx, resp, err)
		return nil, err
	}
	return resp, nil
}

func (c *client) handleCircuitBreakerError(ctx context.Context, err error) (*resty.Response, error) {
	if errors.Is(err, gobreaker.ErrOpenState) {
		c.logCircuitBreakerOpen(ctx, err)
		return nil, fmt.Errorf("circuit breaker open: %w", err)
	}
	return nil, err
}

func (c *client) validateAndReturnResponse(resp *resty.Response) (*resty.Response, error) {
	if err := validateResponse(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) logRequestFailure(ctx context.Context, err error) {
	if c.logging {
		c.logger.Warn(ctx, "failed", map[string]interface{}{
			"event": "request_failed",
			"error": err,
		})
	}
}

func (c *client) logHttpError(ctx context.Context, resp *resty.Response, err error) {
	if c.logging {
		c.logger.Warn(ctx, "http_error", map[string]interface{}{
			"event":  "http_error",
			"status": resp.StatusCode(),
			"error":  err,
		})
	}
}

func (c *client) logCircuitBreakerOpen(ctx context.Context, err error) {
	if c.logging {
		c.logger.Error(ctx, err, "circuit breaker is open", map[string]interface{}{
			"event": "circuit_breaker_open",
			"error": err,
		})
	}
}

func (c *client) Get(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetContext(ctx).Get(c.baseURL + endpoint)
	})
}

func (c *client) Post(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetBody(body).SetContext(ctx).Post(c.baseURL + endpoint)
	})
}

func (c *client) Put(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetBody(body).SetContext(ctx).Put(c.baseURL + endpoint)
	})
}

func (c *client) Patch(ctx context.Context, endpoint string, body interface{}) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetBody(body).SetContext(ctx).Patch(c.baseURL + endpoint)
	})
}

func (c *client) Delete(ctx context.Context, endpoint string) (*resty.Response, error) {
	return c.executeRequest(ctx, func(ctx context.Context) (*resty.Response, error) {
		return c.requester.httpClient.R().SetContext(ctx).Delete(c.baseURL + endpoint)
	})
}

func validateResponse(resp *resty.Response) error {
	if resp == nil {
		return errors.New("response is nil")
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() <= 299 {
		return nil
	}
	bodyPreview := ""
	if resp.Body() != nil && len(resp.Body()) > 0 {
		text := string(resp.Body())
		if len(text) > 200 {
			text = text[:200] + "..."
		}
		bodyPreview = text
	}
	return fmt.Errorf("HTTP %d: %s - %s", resp.StatusCode(), resp.Status(), bodyPreview)
}
