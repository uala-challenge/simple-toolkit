package viper

import (
	"sync"

	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	"github.com/uala-challenge/simple-toolkit/pkg/client/redis"

	"github.com/uala-challenge/simple-toolkit/pkg/client/sns"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Apply() (Config, error)
}

type Config struct {
	Router         simple_router.Config   `json:"router"`
	Log            log.Config             `json:"log"`
	AwsRegion      string                 `json:"aws_region"`
	SQSConfig      *sqs.Config            `json:"sqs"`
	SNSConfig      *sns.Config            `json:"sns"`
	DynamoDBConfig *dynamo.Config         `json:"dynamoDB"`
	RedisConfig    *redis.Config          `json:"redis"`
	Repositories   map[string]interface{} `json:"repositories"`
	UsesCases      map[string]interface{} `json:"usesCases"`
	Endpoints      map[string]interface{} `json:"endpoints"`
}

type service struct {
	propertyFiles []string
	path          string
}

var (
	once     sync.Once
	instance *service
)
