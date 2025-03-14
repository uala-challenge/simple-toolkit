package app_engine

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"

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
	awsCfg := loadAWSConfig(c.Aws, l)
	return &Engine{
		App:                simple_router.NewService(c.Router, l),
		Log:                l,
		SQSClient:          createSQSService(awsCfg, c.SQS),
		SNSClient:          createSNSClient(awsCfg, c.SNS),
		DynamoDBClient:     createDynamoClient(awsCfg, c.Dynamo),
		RedisClient:        createRedisService(c.Redis, l),
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
	return sns2.NewClient(acf, cfg.Endpoint)
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

func loadAWSConfig(ar viper.AwsConfig, l log.Service) aws.Config {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(ar.Region),
	)
	if err != nil {
		l.FatalError(context.Background(), err, map[string]interface{}{})
	}
	return awsCfg
}

func GetConfig[O any](c map[string]interface{}) O {
	h := *new(O)
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   &h,
		TagName:  "json",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	err := decoder.Decode(c)
	if err != nil {
		panic(fmt.Errorf("error loading configuration - cast: %w", err))
	}
	return h
}
