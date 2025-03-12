package app_engine

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
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
	awsCfg := loadAWSConfig(c.AwsConfig, l)
	return &Engine{
		App:                simple_router.NewService(c.Router, l),
		Log:                l,
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
func createSQSService(acf aws.Config, cfg sqs.Config, logger log.Service) sqs.Service[any] {
	return sqs.NewService(acf, cfg, logger)
}

func loadAWSConfig(c aws.Config, l log.Service) aws.Config {
	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(c.Region),
	)
	if err != nil {
		l.FatalError(context.Background(), err, map[string]interface{}{})
	}
	return awsCfg
}
