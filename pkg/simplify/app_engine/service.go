package app_engine

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	rd "github.com/uala-challenge/simple-toolkit/pkg/client/redis"
	sns2 "github.com/uala-challenge/simple-toolkit/pkg/client/sns"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	sqs2 "github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/config/viper"
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
	return &Engine{
		App:                simple_router.NewService(c.Router, l),
		Log:                l,
		SQSService:         createSQSService(awsCfg, c.SQSConfig),
		SNSService:         createSNSClient(awsCfg, c.SNSConfig),
		DynamoDBService:    createDynamoClient(awsCfg, c.DynamoDBConfig),
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
	}, logrus.New())
}
func createSQSService(acf aws.Config, cfg *sqs2.Config) *sqs.Client {
	if cfg == nil {
		return nil
	}
	return sqs2.NewClient(acf, *cfg)
}

func createSNSClient(acf aws.Config, cfg *sns2.Config) *sns.Client {
	if cfg == nil {
		return nil
	}
	return sns2.NewClient(acf, cfg.BaseEndpoint)
}

func createDynamoClient(acf aws.Config, cfg *dynamo.Config) *dynamodb.Client {
	if cfg == nil {
		return nil
	}
	client := dynamo.NewClient(acf, *cfg)
	return client
}

func createRedisService(cfg *rd.Config, l log.Service) *redis.Client {
	if cfg == nil {
		return nil
	}
	client, err := rd.NewClient(*cfg)
	if err != nil {
		l.Error(context.Background(), err, "Error al crear cliente redis", map[string]interface{}{})
	}
	return client
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
