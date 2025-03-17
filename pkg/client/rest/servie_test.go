package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	mockConfigWithRetry = Config{
		BaseURL:          "",
		WithRetry:        true,
		RetryCount:       2,
		RetryWaitTime:    1000 * time.Millisecond,
		RetryMaxWaitTime: 1 * time.Second,
		WithCB:           false,
	}

	mockConfigWithCB = Config{
		BaseURL:            "",
		WithRetry:          false,
		WithCB:             true,
		CBName:             "test_cb",
		CBMaxRequests:      2,
		CBInterval:         5 * time.Second,
		CBTimeout:          3 * time.Second,
		CBRequestThreshold: 4,
		CBFailureRateLimit: 0.5,
	}

	mockHeaders = map[string]string{
		"Authorization": "Bearer test-token",
		"Content-Type":  "application/json",
	}
)

func TestNewClient(t *testing.T) {
	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)
	assert.NotNil(t, client)
}

func TestGetRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)
	resp, err := client.Get(context.Background(), ts.URL, mockHeaders)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPostRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"status": "created"}`))
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)

	body := map[string]string{"name": "test"}
	resp, err := client.Post(context.Background(), ts.URL, body, mockHeaders)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPutRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "updated"}`))
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)

	body := map[string]string{"update": "true"}
	resp, err := client.Put(context.Background(), ts.URL, body, mockHeaders)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestPatchRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "patched"}`))
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)

	body := map[string]string{"patch": "true"}
	resp, err := client.Patch(context.Background(), ts.URL, body, mockHeaders)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteRequest(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)

	resp, err := client.Delete(context.Background(), ts.URL, mockHeaders)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestRetryMechanism(t *testing.T) {
	attempts := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			http.Error(w, "Temporary error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)

	resp, err := client.Get(context.Background(), ts.URL, mockHeaders)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestCircuitBreaker(t *testing.T) {
	failures := 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		failures++
		if failures < 3 {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "success"}`))
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithCB, l)

	// Realizar varias solicitudes para simular fallos y recuperación
	var resp *resty.Response
	var err error

	for i := 0; i < 5; i++ {
		resp, err = client.Get(context.Background(), ts.URL, mockHeaders)
		time.Sleep(500 * time.Millisecond) // Esperar para que el breaker se reajuste
	}

	assert.NoError(t, err, "Circuit Breaker debería permitir la última solicitud exitosa")
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestExponentialBackoffWithJitter(t *testing.T) {
	l := logrus.New()
	initialWait := 100 * time.Millisecond
	maxWait := 5 * time.Second

	waitTime := exponentialBackoffWithJitter(initialWait, maxWait, 3, l)

	assert.GreaterOrEqual(t, waitTime, initialWait)
	assert.LessOrEqual(t, waitTime, maxWait)
}

func TestSetDefaultConfig(t *testing.T) {
	cfg := Config{} // Sin valores definidos

	setDefaultConfig(&cfg)

	assert.Equal(t, DefaultRetryCount, cfg.RetryCount)
	assert.Equal(t, DefaultRetryWaitTime, cfg.RetryWaitTime)
	assert.Equal(t, DefaultRetryMaxWaitTime, cfg.RetryMaxWaitTime)
	assert.Equal(t, DefaultCBMaxRequests, cfg.CBMaxRequests)
	assert.Equal(t, DefaultCBInterval, cfg.CBInterval)
	assert.Equal(t, DefaultCBTimeout, cfg.CBTimeout)
	assert.Equal(t, DefaultCBRequestThreshold, cfg.CBRequestThreshold)
	assert.Equal(t, DefaultCBFailureRateLimit, cfg.CBFailureRateLimit)
}
