package app_engine

import (
	"context"
	"fmt"
	"os"

	"go.elastic.co/ecslogrus"

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
	tracer := logrus.New()
	tracer.SetOutput(os.Stdout)
	tracer.SetFormatter(&ecslogrus.Formatter{})
	tracer.Level = logrus.DebugLevel
	v := viper.NewService(tracer)
	c, err := v.Apply()
	if err != nil {
		tracer.Fatal(err)
	}
	awsCfg := loadAWSConfig(c.Aws, tracer)
	return &Engine{
		App:                simple_router.NewService(c.Router),
		SQSClient:          createSQSService(awsCfg, c.SQS, tracer),
		SNSClient:          createSNSClient(awsCfg, c.SNS, tracer),
		DynamoDBClient:     createDynamoClient(awsCfg, c.Dynamo, tracer),
		RedisClient:        createRedisService(c.Redis, tracer),
		RepositoriesConfig: c.Repositories,
		UsesCasesConfig:    c.Cases,
		HandlerConfig:      c.Endpoints,
		Log:                configLogLevel(c.Log, tracer),
	}
}

func configLogLevel(c log.Config, l *logrus.Logger) log.Service {
	return log.NewService(log.Config{
		Level: c.Level,
		Path:  c.Path,
	}, l)
}
func createSQSService(acf aws.Config, cfg *sqs2.Config, l *logrus.Logger) *sqs.Client {
	if cfg == nil {
		return nil
	}
	return sqs2.NewClient(acf, *cfg, l)
}

func createSNSClient(acf aws.Config, cfg *sns2.Config, l *logrus.Logger) *sns.Client {
	if cfg == nil {
		return nil
	}
	return sns2.NewClient(acf, cfg.Endpoint, l)
}

func createDynamoClient(acf aws.Config, cfg *dynamo.Config, l *logrus.Logger) *dynamodb.Client {
	if cfg == nil {
		return nil
	}
	client := dynamo.NewClient(acf, *cfg, l)
	return client
}

func createRedisService(cfg *rd.Config, l *logrus.Logger) *redis.Client {
	if cfg == nil {
		return nil
	}
	client, err := rd.NewClient(*cfg, l)
	if err != nil {
		l.Error(err)
	}
	return client
}

func loadAWSConfig(ar viper.AwsConfig, l *logrus.Logger) aws.Config {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(ar.Region),
	)
	if err != nil {
		l.Fatal(err)
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
