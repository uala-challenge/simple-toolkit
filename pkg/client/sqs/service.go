package sqs

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"go.opentelemetry.io/otel/trace"
)

type sqsClient interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}
type service struct {
	client      sqsClient
	config      Config
	logger      log.Service
	tracer      trace.Tracer
	retryDelays []time.Duration
}

var _ Service = (*service)(nil)

func NewService(acf aws.Config, cfg Config, logger log.Service, tracer trace.Tracer) Service {
	client := createSQSClient(acf, cfg, logger)

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

func createSQSClient(acf aws.Config, cfg Config, logger log.Service) sqsClient {
	return sqs.NewFromConfig(acf, func(o *sqs.Options) {
		if cfg.BaseEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.BaseEndpoint)
			logger.Info(context.TODO(), "Configurando SQS con LocalStack", map[string]interface{}{
				"endpoint": cfg.BaseEndpoint,
			})
		} else {
			logger.Info(context.TODO(), "Configurando SQS con AWS", nil)
		}
	})
}

func (s *service) ReceiveMessage(ctx context.Context) (Message, error) {
	ctx, span := s.tracer.Start(ctx, "Receive SQS Message")
	defer span.End()

	span.SetAttributes(attribute.String("sqs.queue_url", s.config.QueueURL))

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
	message := Message{
		ID:            *msg.MessageId,
		Body:          *msg.Body,
		ReceiptHandle: *msg.ReceiptHandle,
	}

	span.SetAttributes(attribute.String("sqs.message_id", message.ID))

	return message, nil
}

func (s *service) DeleteMessage(ctx context.Context, receiptHandle string) error {
	ctx, span := s.tracer.Start(ctx, "Delete SQS Message")
	defer span.End()

	span.SetAttributes(attribute.String("sqs.queue_url", s.config.QueueURL))
	span.SetAttributes(attribute.String("sqs.receipt_handle", receiptHandle))

	var lastErr error
	for attempt, delay := range s.retryDelays {
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
			"attempt":       attempt + 1,
			"error":         err.Error(),
			"receiptHandle": receiptHandle,
		})

		time.Sleep(delay)
	}

	s.logger.Error(ctx, lastErr, "No se pudo eliminar el mensaje después de intentos máximos", map[string]interface{}{
		"receiptHandle": receiptHandle,
	})
	return lastErr
}
