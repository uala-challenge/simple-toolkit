package sqs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type MockSQSClient struct {
	mock.Mock
}

func (m *MockSQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

func (m *MockSQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*sqs.DeleteMessageOutput), args.Error(1)
}

func TestCreateClienteWithError(t *testing.T) {
	mockClient := new(MockSQSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sns-test")

	awsCfg := aws.Config{
		Region: "us-east-1",
	}

	cfg := Config{
		QueueURL:        "https://sqs.us-east-1.amazonaws.com/000000000000/test",
		BaseEndpoint:    "",
		MaxRetries:      2,
		MaxMessages:     1,
		WaitTimeSeconds: 1,
		TimeoutSeconds:  2,
	}

	service := NewService(awsCfg, cfg, mockLogger, mockTracer)

	mockClient.On("ReceiveMessage", mock.Anything, &sqs.ReceiveMessageInput{
		QueueUrl:            &cfg.QueueURL,
		MaxNumberOfMessages: cfg.MaxMessages,
		WaitTimeSeconds:     cfg.WaitTimeSeconds,
	}).Return(&sqs.ReceiveMessageOutput{Messages: []types.Message{}}, nil).Once()

	_, err := service.ReceiveMessage(context.Background())

	assert.Error(t, err)
}

func TestCreateClienteWithError2(t *testing.T) {
	mockClient := new(MockSQSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sns-test")

	awsCfg := aws.Config{
		Region: "us-east-1",
	}

	cfg := Config{
		QueueURL:        "https://sqs.us-east-1.amazonaws.com/000000000000/test",
		BaseEndpoint:    "test",
		MaxRetries:      2,
		MaxMessages:     1,
		WaitTimeSeconds: 1,
		TimeoutSeconds:  2,
	}

	service := NewService(awsCfg, cfg, mockLogger, mockTracer)

	mockClient.On("ReceiveMessage", mock.Anything, &sqs.ReceiveMessageInput{
		QueueUrl:            &cfg.QueueURL,
		MaxNumberOfMessages: cfg.MaxMessages,
		WaitTimeSeconds:     cfg.WaitTimeSeconds,
	}).Return(&sqs.ReceiveMessageOutput{Messages: []types.Message{}}, nil).Once()

	_, err := service.ReceiveMessage(context.Background())

	assert.Error(t, err)
}

func TestReceiveMessage_Success(t *testing.T) {
	mockClient := new(MockSQSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sqs-test")

	cfg := Config{
		QueueURL:        "https://sqs.us-east-1.amazonaws.com/000000000000/test",
		BaseEndpoint:    "",
		MaxRetries:      2,
		MaxMessages:     1,
		WaitTimeSeconds: 1,
		TimeoutSeconds:  2,
	}

	service := newServiceWithClient(mockClient, cfg, mockLogger, mockTracer)

	mockClient.On("ReceiveMessage", mock.Anything, &sqs.ReceiveMessageInput{
		QueueUrl:            &cfg.QueueURL,
		MaxNumberOfMessages: cfg.MaxMessages,
		WaitTimeSeconds:     cfg.WaitTimeSeconds,
	}).Return(&sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{
				Body:          aws.String("Hello, SQS!"),
				MessageId:     aws.String("msg-1"),
				ReceiptHandle: aws.String("handle-1"),
			},
		},
	}, nil).Once()

	msg, err := service.ReceiveMessage(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "msg-1", msg.ID)
	assert.Equal(t, "Hello, SQS!", msg.Body)
	assert.Equal(t, "handle-1", msg.ReceiptHandle)
	mockClient.AssertExpectations(t)
}

func TestReceiveMessage_NoMessages(t *testing.T) {
	mockClient := new(MockSQSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sqs-test")

	cfg := Config{
		QueueURL:        "https://sqs.us-east-1.amazonaws.com/000000000000/test",
		BaseEndpoint:    "",
		MaxRetries:      2,
		MaxMessages:     1,
		WaitTimeSeconds: 1,
		TimeoutSeconds:  2,
	}

	service := newServiceWithClient(mockClient, cfg, mockLogger, mockTracer)

	mockClient.On("ReceiveMessage", mock.Anything, &sqs.ReceiveMessageInput{
		QueueUrl:            &cfg.QueueURL,
		MaxNumberOfMessages: cfg.MaxMessages,
		WaitTimeSeconds:     cfg.WaitTimeSeconds,
	}).Return(&sqs.ReceiveMessageOutput{Messages: []types.Message{}}, nil).Once()

	msg, err := service.ReceiveMessage(context.Background())

	assert.ErrorIs(t, err, ErrNoMessages)
	assert.Empty(t, msg)
	mockClient.AssertExpectations(t)
}

func TestReceiveMessage_Error(t *testing.T) {
	mockClient := new(MockSQSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sqs-test")

	cfg := Config{
		QueueURL:        "https://sqs.us-east-1.amazonaws.com/000000000000/test",
		BaseEndpoint:    "",
		MaxRetries:      2,
		MaxMessages:     1,
		WaitTimeSeconds: 1,
		TimeoutSeconds:  2,
	}

	service := newServiceWithClient(mockClient, cfg, mockLogger, mockTracer)

	mockClient.On("ReceiveMessage", mock.Anything, mock.Anything).
		Return(&sqs.ReceiveMessageOutput{}, errors.New("AWS SQS Error")).Once()

	msg, err := service.ReceiveMessage(context.Background())

	assert.Error(t, err)
	assert.Equal(t, "AWS SQS Error", err.Error())
	assert.Empty(t, msg)
	mockClient.AssertExpectations(t)
}

func TestDeleteMessage_Success(t *testing.T) {
	mockClient := new(MockSQSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sqs-test")

	cfg := Config{
		QueueURL:        "https://sqs.us-east-1.amazonaws.com/000000000000/test",
		BaseEndpoint:    "",
		MaxRetries:      2,
		MaxMessages:     1,
		WaitTimeSeconds: 1,
		TimeoutSeconds:  2,
	}

	service := newServiceWithClient(mockClient, cfg, mockLogger, mockTracer)

	mockClient.On("DeleteMessage", mock.Anything, &sqs.DeleteMessageInput{
		QueueUrl:      &cfg.QueueURL,
		ReceiptHandle: aws.String("handle-1"),
	}).Return(&sqs.DeleteMessageOutput{}, nil).Once()

	err := service.DeleteMessage(context.Background(), "handle-1")

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestDeleteMessage_Failure(t *testing.T) {
	mockClient := new(MockSQSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sqs-test")

	cfg := Config{
		QueueURL:        "https://sqs.us-east-1.amazonaws.com/000000000000/test",
		BaseEndpoint:    "",
		MaxRetries:      2,
		MaxMessages:     1,
		WaitTimeSeconds: 2,
		TimeoutSeconds:  2,
	}

	service := newServiceWithClient(mockClient, cfg, mockLogger, mockTracer)

	mockClient.On("DeleteMessage", mock.Anything, &sqs.DeleteMessageInput{
		QueueUrl:      &cfg.QueueURL,
		ReceiptHandle: aws.String("handle-1"),
	}).Return(&sqs.DeleteMessageOutput{}, errors.New("AWS SQS Delete Error")).Times(cfg.MaxRetries)

	err := service.DeleteMessage(context.Background(), "handle-1")

	assert.Error(t, err)
	assert.Equal(t, "AWS SQS Delete Error", err.Error())
	mockClient.AssertExpectations(t)
}

func newServiceWithClient(client sqsClient, cfg Config, logger log.Service, tracer trace.Tracer) Service {
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
