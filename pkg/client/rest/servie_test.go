package rest

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker/v2"
	"github.com/stretchr/testify/assert"
)

var (
	mockConfigWithRetry = Config{
		BaseURL:            "",
		WithRetry:          true,
		RetryCount:         2,
		RetryWaitTime:      1000 * time.Millisecond,
		RetryMaxWaitTime:   1 * time.Second,
		WithCB:             false,
		CBName:             "",
		CBMaxRequests:      0,
		CBInterval:         0,
		CBTimeout:          0,
		CBRequestThreshold: 0,
		CBFailureRateLimit: 0,
	}
	mockConfigWithCB = Config{
		BaseURL:            "",
		WithRetry:          false,
		RetryCount:         0,
		RetryWaitTime:      0,
		RetryMaxWaitTime:   0,
		WithCB:             true,
		CBName:             "test_cb",
		CBMaxRequests:      2,
		CBInterval:         5 * time.Second,
		CBTimeout:          3 * time.Second,
		CBRequestThreshold: 4,
		CBFailureRateLimit: 0.5,
	}
	mockConfigWithCBPort = Config{
		BaseURL:            "http://localhost:9999",
		WithRetry:          false,
		RetryCount:         0,
		RetryWaitTime:      0,
		RetryMaxWaitTime:   0,
		WithCB:             true,
		CBName:             "test_cb",
		CBMaxRequests:      2,
		CBInterval:         5 * time.Second,
		CBTimeout:          3 * time.Second,
		CBRequestThreshold: 4,
		CBFailureRateLimit: 0.5,
	}
)

func TestNewClient(t *testing.T) {
	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)
	assert.NotNil(t, client)
}

func TestGetRequest(t *testing.T) {
	attempts := 0

	_ = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(nil)
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(nil)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)
	_, err := client.Get(context.Background(), ts.URL)

	assert.Error(t, err)
}

func TestGetRequestWithError(t *testing.T) {
	failures := 0
	maxFailures := 2

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if failures < maxFailures {
			failures++
			http.Error(w, "simulated server error", http.StatusInternalServerError)

			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"status": "created"}`))
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithCBPort, l)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test"}`))
	var err error

	for i := 0; i < maxFailures+1; i++ {
		_, err = client.Post(context.Background(), ts.URL, requestBody)
	}

	assert.Error(t, err)
}

func TestPostRequestWithCB(t *testing.T) {
	failures := 0
	maxFailures := 3

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if failures < 2 {
			failures++
			http.Error(w, "simulated server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"status": "created"}`))
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithCB, l)
	requestBody := bytes.NewBuffer([]byte(`{"name": "test"}`))
	var err error
	var response *resty.Response

	for i := 0; i < maxFailures+1; i++ {
		response, err = client.Post(context.Background(), ts.URL, requestBody)
	}

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPostRequestWithCBAndError(t *testing.T) {
	failures := 0

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if failures < 4 {
			failures++
			http.Error(w, "simulated server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"status": "created"}`))
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithCB, l)
	client.WithLogging(true)
	requestBody := bytes.NewBuffer([]byte(`{"name": "test"}`))
	response, err := client.Post(context.Background(), ts.URL, requestBody)

	assert.Error(t, err)
	assert.Nil(t, response)
}

func TestPostRequestIsOpen(t *testing.T) {
	failures := 0
	maxFailures := 3

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if failures < maxFailures {
			failures++
			http.Error(w, "simulated server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, err := w.Write([]byte(`{"status": "created"}`))
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	l := logrus.New()
	client := NewClient(mockConfigWithCB, l)

	requestBody := bytes.NewBuffer([]byte(`{"name": "test"}`))
	var err error
	var response *resty.Response

	for i := 0; i < maxFailures+1; i++ {
		response, err = client.Post(context.Background(), ts.URL, requestBody)
	}

	assert.Error(t, err, "circuit breaker is open")
	assert.Nil(t, response)
}

func TestPutRequest(t *testing.T) {
	attempts := 0

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"message": "success"}`))
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()
	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)
	var err error
	var response *resty.Response

	for i := 0; i < 3; i++ {
		response, err = client.Put(context.Background(), ts.URL, nil)
	}
	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestPatchRequest(t *testing.T) {
	mockResponse := `{"message": "success"}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PATCH", r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockResponse))
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()
	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)
	response, err := client.Patch(context.Background(), ts.URL, nil)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestDeleteRequest(t *testing.T) {
	mockResponse := `{"message": "success"}`
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(mockResponse))
		if err != nil {
			return
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()
	l := logrus.New()
	client := NewClient(mockConfigWithRetry, l)
	response, err := client.Delete(context.Background(), ts.URL)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestExponentialBackoffWithJitter(t *testing.T) {
	l := logrus.New()
	initialWait := 100 * time.Millisecond
	maxWait := 5 * time.Second

	waitTimeZero := exponentialBackoffWithJitter(initialWait, maxWait, 0, l)
	expectedBaseZero := initialWait * time.Duration(math.Pow(BackoffFactor, 0))
	assert.GreaterOrEqual(t, waitTimeZero, expectedBaseZero, "Backoff incorrecto para attempt=0")

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

func TestCheckBreakerState(t *testing.T) {
	l := logrus.New()
	cfg := Config{CBMaxRequests: 3, CBRequestThreshold: 5, CBFailureRateLimit: 0.5}

	counts := gobreaker.Counts{
		Requests:            10,
		TotalFailures:       6,
		ConsecutiveFailures: 4,
	}
	result := checkBreakerState(counts, cfg, l)
	assert.True(t, result)
}
func TestRetryAfterFunc(t *testing.T) {
	l := logrus.New()
	client := resty.New()
	resp := &resty.Response{Request: &resty.Request{Attempt: 3}}
	retryWait, err := retryAfterFunc(DefaultRetryWaitTime, DefaultRetryMaxWaitTime, l)(client, resp)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, retryWait, DefaultRetryWaitTime)
	assert.LessOrEqual(t, retryWait, DefaultRetryMaxWaitTime)
}

func TestValidateResponse_Success(t *testing.T) {
	resp := &resty.Response{RawResponse: &http.Response{StatusCode: http.StatusOK}}
	err := validateResponse(resp)

	assert.NoError(t, err)
}
