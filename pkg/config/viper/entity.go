package viper

import (
	"github.com/uala-challenge/simple-toolkit/pkg/client/sns"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/database/dynamo"
	"github.com/uala-challenge/simple-toolkit/pkg/database/redis"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"sync"
)

type Service interface {
	Apply() (Config, error)
}

type Config struct {
	Router         simple_router.Config   `json:"router"`
	Log            log.Config             `json:"log"`
	AwsRegion      string                 `json:"awsRegion"`
	SQSConfig      *sqs.Config            `json:"sqsConfig"`
	SNSConfig      *sns.Config            `json:"snsConfig"`
	DynamoDBConfig *dynamo.Config         `json:"dynamoDBConfig"`
	RedisConfig    *redis.Config          `json:"redisConfig"`
	Repositories   map[string]interface{} `json:"repositories"`
	UsesCases      map[string]interface{} `json:"uses_cases"`
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
