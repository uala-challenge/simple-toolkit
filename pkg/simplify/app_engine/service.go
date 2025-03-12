package app_engine

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sns"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/config/viper"
	"github.com/uala-challenge/simple-toolkit/pkg/database/dynamo"
	"github.com/uala-challenge/simple-toolkit/pkg/database/redis"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

func NewApp() *Engine {
	v := viper.NewService()
	c, err := v.Apply()
	if err != nil {
		panic(err)
	}
	l := creteLog(c.Log)
	awsCfg := loadAWSConfig(c.AwsRegion, l)
	tracerProvider := otel.GetTracerProvider()
	return &Engine{
		App:                simple_router.NewService(c.Router, l),
		Log:                l,
		SQSService:         createSQSService(awsCfg, c.SQSConfig, l),
		SNSService:         createSNSService(awsCfg, c.SNSConfig, l, tracerProvider.Tracer("sns-service")),
		DynamoDBService:    createDynamoService(awsCfg, c.DynamoDBConfig, l),
		RedisService:       createRedisService(c.RedisConfig, l),
		RepositoriesConfig: c.Repositories,
		UsesCasesConfig:    c.UsesCases,
		HandlerConfig:      c.Endpoints,
	}
}

func creteLog(c log.Config) log.Service {
	return log.NewService(log.Config{
		Level: c.Level,
		Path:  c.Path,
	})
}
func createSQSService(acf aws.Config, cfg *sqs.Config, logger log.Service) sqs.Service {
	if cfg == nil {
		return nil
	}
	return sqs.NewService(acf, *cfg, logger)
}

func createSNSService(acf aws.Config, cfg *sns.Config, logger log.Service, tracer trace.Tracer) sns.Service {
	if cfg == nil {
		return nil
	}
	return sns.NewService(acf, *cfg, logger, tracer)
}

func createDynamoService(acf aws.Config, cfg *dynamo.Config, logger log.Service) dynamo.Service {
	if cfg == nil {
		return nil
	}
	return dynamo.NewService(acf, *cfg, logger)
}

func createRedisService(cfg *redis.Config, logger log.Service) redis.Service {
	if cfg == nil {
		return nil
	}
	return redis.NewService(*cfg, logger)
}

func loadAWSConfig(ar string, l log.Service) aws.Config {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(ar),
	)
	if err != nil {
		l.FatalError(context.Background(), err, map[string]interface{}{})
	}
	return awsCfg
}
