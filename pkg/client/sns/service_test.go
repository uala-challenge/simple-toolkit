package sns

import (
	"context"
	"encoding/json"
	"go.opentelemetry.io/otel"
	"testing"

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

func TestPublishMessage_Success(t *testing.T) {
	mockClient := new(MockSNSClient)
	mockLogger := log.NewService(log.Config{Level: "info"})
	mockTracer := otel.GetTracerProvider().Tracer("sns-test") // ✅ Solución

	awsCfg := aws.Config{
		Region: "us-east-1",
	}

	cfg := Config{
		TopicARN:       "arn:aws:sns:us-east-1:000000000000:uala-challenge",
		BaseEndpoint:   "http://localhost:4566",
		MaxRetries:     3,
		TimeoutSeconds: 5,
	}

	service := NewService(awsCfg, cfg, mockLogger, mockTracer)

	message := map[string]string{"key": "value"}
	messageJSON, _ := json.Marshal(message)

	mockClient.On("Publish", mock.Anything, &sns.PublishInput{
		TopicArn: aws.String(cfg.TopicARN),
		Message:  aws.String(string(messageJSON)),
	}).Return(&sns.PublishOutput{MessageId: aws.String("12345")}, nil)

	err := service.PublishMessage(context.Background(), message)

	assert.NoError(t, err)
}
