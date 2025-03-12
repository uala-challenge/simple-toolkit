package sns

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type MockSNSClient struct {
	mock.Mock
}

func (m *MockSNSClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sns.PublishOutput), args.Error(1)
}

func TestCriateClienteWithErrorPublish(t *testing.T) {
	mockClient := new(MockSNSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sns-test")

	awsCfg := aws.Config{
		Region: "us-east-1",
	}

	cfg := Config{
		TopicARN:       "arn:aws:sns:us-east-1:000000000000:test",
		BaseEndpoint:   "",
		MaxRetries:     2,
		TimeoutSeconds: 2,
	}

	service := NewService(awsCfg, cfg, mockLogger, mockTracer)

	message := map[string]string{"key": "value"}
	messageJSON, _ := json.Marshal(message)

	mockClient.On("Publish", mock.Anything, &sns.PublishInput{
		TopicArn: aws.String(cfg.TopicARN),
		Message:  aws.String(string(messageJSON)),
	}).Return(&sns.PublishOutput{MessageId: aws.String("12345")}, nil)

	err := service.PublishMessage(context.Background(), message)

	assert.Error(t, err)
}

func TestCriateClienteWithError(t *testing.T) {
	mockClient := new(MockSNSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sns-test") // ✅ Solución

	awsCfg := aws.Config{
		Region: "us-east-1",
	}

	cfg := Config{
		TopicARN:       "arn:aws:sns:us-east-1:000000000000:test",
		BaseEndpoint:   "test",
		MaxRetries:     2,
		TimeoutSeconds: 2,
	}

	service := NewService(awsCfg, cfg, mockLogger, mockTracer)

	message := map[string]string{"key": "value"}
	messageJSON, _ := json.Marshal(message)

	mockClient.On("Publish", mock.Anything, &sns.PublishInput{
		TopicArn: aws.String(cfg.TopicARN),
		Message:  aws.String(string(messageJSON)),
	}).Return(&sns.PublishOutput{MessageId: aws.String("12345")}, nil)

	err := service.PublishMessage(context.Background(), message)

	assert.Error(t, err)
}

func TestPublishMessage_Success(t *testing.T) {
	mockClient := new(MockSNSClient) // ✅ Creamos el mock
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sns-test")

	cfg := Config{
		TopicARN:       "arn:aws:sns:us-east-1:000000000000:uala-challenge",
		BaseEndpoint:   "",
		MaxRetries:     2,
		TimeoutSeconds: 2,
	}

	service := newServiceWithClient(mockClient, cfg, mockLogger, mockTracer)

	message := map[string]string{"key": "value"}
	messageJSON, _ := json.Marshal(message)

	mockClient.On("Publish", mock.Anything, &sns.PublishInput{
		TopicArn: aws.String(cfg.TopicARN),
		Message:  aws.String(string(messageJSON)),
	}).Return(&sns.PublishOutput{MessageId: aws.String("12345")}, nil).Once()

	err := service.PublishMessage(context.Background(), message)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func newServiceWithClient(client snsClient, cfg Config, logger log.Service, tracer trace.Tracer) Service {
	retryDelays := make([]time.Duration, cfg.MaxRetries)
	for i := 0; i < cfg.MaxRetries; i++ {
		retryDelays[i] = time.Duration(1<<i) * time.Second
	}

	return &service{
		client:      client,
		config:      cfg,
		logger:      logger,
		tracer:      tracer,
		retryDelays: retryDelays,
	}
}

func TestPublishMessage_JSONMarshalError(t *testing.T) {
	mockClient := new(MockSNSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sns-test")

	cfg := Config{
		TopicARN:       "arn:aws:sns:us-east-1:000000000000:uala-challenge",
		BaseEndpoint:   "",
		MaxRetries:     2,
		TimeoutSeconds: 2,
	}

	service := newServiceWithClient(mockClient, cfg, mockLogger, mockTracer)

	invalidMessage := make(chan int)

	err := service.PublishMessage(context.Background(), invalidMessage)

	assert.Error(t, err)
	mockClient.AssertNotCalled(t, "Publish")
}
