package app_engine

import (
	"github.com/redis/go-redis/v9"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sns"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Engine struct {
	App                *simple_router.App
	Log                log.Service
	SQSClient          *sqs.Sqs
	SNSClient          *sns.Sns
	DynamoDBClient     *dynamo.Dynamo
	RedisClient        *redis.Client
	RepositoriesConfig map[string]interface{}
	UsesCasesConfig    map[string]interface{}
	HandlerConfig      map[string]interface{}
}
