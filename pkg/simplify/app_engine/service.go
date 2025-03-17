package app_engine

import (
	"context"
	"fmt"
	"github.com/uala-challenge/simple-toolkit/pkg/client/rest"
	"os"

	"go.elastic.co/ecslogrus"

	"github.com/mitchellh/mapstructure"

	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	rd "github.com/uala-challenge/simple-toolkit/pkg/client/redis"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sns"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
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
		RestClients:        createHttpClient(c.Rest, tracer),
		Log:                configLogLevel(c.Log, tracer),
	}
}

func configLogLevel(c log.Config, l *logrus.Logger) log.Service {
	return log.NewService(log.Config{
		Level: c.Level,
		Path:  c.Path,
	}, l)
}
func createSQSService(acf aws.Config, cfg *sqs.Config, l *logrus.Logger) *sqs.Sqs {
	if cfg == nil {
		return nil
	}
	return sqs.NewClient(acf, *cfg, l)
}

func createSNSClient(acf aws.Config, cfg *sns.Config, l *logrus.Logger) *sns.Sns {
	if cfg == nil {
		return nil
	}
	return sns.NewClient(acf, cfg.Endpoint, l)
}

func createDynamoClient(acf aws.Config, cfg *dynamo.Config, l *logrus.Logger) *dynamo.Dynamo {
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

func createHttpClient(c []map[string]rest.Config, l *logrus.Logger) map[string]rest.Service {
	httpClients := make(map[string]rest.Service)
	for _, v := range c {
		for k, v := range v {
			httpClients[k] = rest.NewClient(v, l)
		}
	}
	return httpClients
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
