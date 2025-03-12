package sqs

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *sqs.Client
	config Config
	logger log.Service
}

var _ Service = (*service)(nil)

func NewService(acf aws.Config, cfg Config, logger log.Service) Service {
	client := sqs.NewFromConfig(acf, func(o *sqs.Options) {
		if cfg.BaseEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.BaseEndpoint)
			logger.Info(context.TODO(), "🔧 Configurando SQS con LocalStack", map[string]interface{}{
				"endpoint": cfg.BaseEndpoint,
			})
		} else {
			logger.Info(context.TODO(), "🚀 Configurando SQS con AWS", nil)
		}
	})

	return &service{
		client: client,
		config: cfg,
		logger: logger,
	}
}

func (s *service) ReceiveMessage(ctx context.Context) (Message, error) {
	s.logger.Info(ctx, "Recibiendo un mensaje de SQS", map[string]interface{}{
		"queue": s.config.QueueURL,
	})

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.TimeoutSeconds)*time.Second)
	defer cancel()

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            &s.config.QueueURL,
		MaxNumberOfMessages: s.config.MaxMessages,
		WaitTimeSeconds:     s.config.WaitTimeSeconds,
	}

	result, err := s.client.ReceiveMessage(ctx, input)
	if err != nil {
		s.logger.Error(ctx, err, "Error al recibir mensaje de SQS", map[string]interface{}{
			"queue": s.config.QueueURL,
		})
		return Message{}, err
	}

	if len(result.Messages) == 0 {
		return Message{}, ErrNoMessages
	}

	msg := result.Messages[0]
	return Message{
		ID:            *msg.MessageId,
		Body:          *msg.Body,
		ReceiptHandle: *msg.ReceiptHandle,
	}, nil
}

func (s *service) DeleteMessage(ctx context.Context, receiptHandle string) error {
	retryDelays := make([]time.Duration, s.config.MaxRetries)
	for i := 0; i < s.config.MaxRetries; i++ {
		retryDelays[i] = time.Duration(1<<i) * time.Second
	}

	var lastErr error
	for attempt := 1; attempt <= len(retryDelays); attempt++ {
		_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      &s.config.QueueURL,
			ReceiptHandle: &receiptHandle,
		})
		if err == nil {
			s.logger.Info(ctx, "Mensaje eliminado correctamente de SQS", nil)
			return nil
		}

		lastErr = err
		s.logger.Warn(ctx, "No se pudo eliminar el mensaje", map[string]interface{}{
			"attempt":       attempt,
			"error":         err.Error(),
			"receiptHandle": receiptHandle,
		})

		time.Sleep(retryDelays[attempt-1])
	}

	s.logger.Error(ctx, lastErr, "No se pudo eliminar el mensaje después de intentos máximos", map[string]interface{}{
		"receiptHandle": receiptHandle,
	})
	return lastErr
}
