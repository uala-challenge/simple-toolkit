package sns

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type snsClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

type service struct {
	client      snsClient
	config      Config
	logger      log.Service
	tracer      trace.Tracer
	retryDelays []time.Duration
}

var _ Service = (*service)(nil)

func NewService(acf aws.Config, cfg Config, logger log.Service, tracer trace.Tracer) Service {
	client := createSNSClient(acf, cfg, logger)

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

func createSNSClient(acf aws.Config, cfg Config, logger log.Service) snsClient {
	return sns.NewFromConfig(acf, func(o *sns.Options) {
		if cfg.BaseEndpoint != "" {
			o.BaseEndpoint = aws.String(cfg.BaseEndpoint)
			logger.Info(context.TODO(), "🔧 Configurando SNS con LocalStack", map[string]interface{}{
				"endpoint": cfg.BaseEndpoint,
			})
		} else {
			logger.Info(context.TODO(), "🚀 Configurando SNS con AWS", nil)
		}
	})
}

func (s *service) PublishMessage(ctx context.Context, message interface{}) error {
	ctx, span := s.tracer.Start(ctx, "Publish SNS Message")
	defer span.End()

	messageID := uuid.New().String()
	span.SetAttributes(attribute.String("sns.message_id", messageID))
	span.SetAttributes(attribute.String("sns.topic_arn", s.config.TopicARN))

	s.logger.Info(ctx, "📢 Publicando mensaje en SNS", map[string]interface{}{
		"message_id": messageID,
		"topic_arn":  s.config.TopicARN,
	})

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.TimeoutSeconds)*time.Second)
	defer cancel()

	messageJSON, err := json.Marshal(message)
	if err != nil {
		s.logger.Error(ctx, err, "Error serializando el mensaje para SNS", nil)
		return err
	}

	var lastErr error
	for attempt := 1; attempt <= s.config.MaxRetries; attempt++ {
		input := &sns.PublishInput{
			TopicArn: &s.config.TopicARN,
			Message:  aws.String(string(messageJSON)),
		}

		result, err := s.client.Publish(ctx, input)
		if err == nil {
			s.logger.Info(ctx, "Mensaje publicado en SNS", map[string]interface{}{
				"message_id": *result.MessageId,
				"topic_arn":  s.config.TopicARN,
			})
			return nil
		}

		lastErr = err
		s.logger.Warn(ctx, "Reintentando publicación de mensaje en SNS...", map[string]interface{}{
			"attempt":    attempt,
			"message_id": messageID,
			"error":      err.Error(),
			"topic_arn":  s.config.TopicARN,
		})

		if attempt < s.config.MaxRetries {
			time.Sleep(s.retryDelays[attempt-1])
		}
	}

	s.logger.Error(ctx, lastErr, "No se pudo publicar el mensaje después de intentos máximos", map[string]interface{}{
		"message_id": messageID,
		"topic_arn":  s.config.TopicARN,
	})
	return lastErr
}
