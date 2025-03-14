package app_engine

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Engine struct {
	App                *simple_router.App
	Log                log.Service
	SQSClient          *sqs.Client
	SNSClient          *sns.Client
	DynamoDBClient     *dynamodb.Client
	RedisClient        *redis.Client
	RepositoriesConfig map[string]interface{}
	UsesCasesConfig    map[string]interface{}
	HandlerConfig      map[string]interface{}
}
