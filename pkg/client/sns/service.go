package sns

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *sns.Client
	config Config
	logger log.Service
}

var _ Service = (*service)(nil)

func NewService(acf aws.Config, cfg Config, logger log.Service) Service {
	client := sns.NewFromConfig(acf)
	return &service{
		client: client,
		config: cfg,
		logger: logger,
	}
}

func (s *service) Accept(ctx context.Context, message interface{}) error {
	s.logger.Info(ctx, "📢 Publicando mensaje en SNS", map[string]interface{}{
		"topic_arn": s.config.TopicARN,
	})

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.config.TimeoutSeconds)*time.Second)
	defer cancel()

	messageJSON, err := json.Marshal(message)
	if err != nil {
		s.logger.Error(ctx, err, "Error serializando el mensaje para SNS", nil)
		return err
	}

	retryDelays := make([]time.Duration, s.config.MaxRetries)
	for i := 0; i < s.config.MaxRetries; i++ {
		retryDelays[i] = time.Duration(1<<i) * time.Second
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
		s.logger.Warn(ctx, map[string]interface{}{
			"attempt":   attempt,
			"error":     err.Error(),
			"topic_arn": s.config.TopicARN,
		})

		if attempt < s.config.MaxRetries {
			time.Sleep(retryDelays[attempt-1])
		}
	}

	s.logger.Error(ctx, lastErr, "No se pudo publicar el mensaje después de intentos máximos", map[string]interface{}{
		"topic_arn":   s.config.TopicARN,
		"max_retries": s.config.MaxRetries,
	})
	return lastErr
}
