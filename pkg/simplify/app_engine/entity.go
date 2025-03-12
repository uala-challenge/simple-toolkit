package app_engine

import (
	"github.com/uala-challenge/simple-toolkit/pkg/client/sns"
	"github.com/uala-challenge/simple-toolkit/pkg/client/sqs"
	"github.com/uala-challenge/simple-toolkit/pkg/database/dynamo"
	"github.com/uala-challenge/simple-toolkit/pkg/database/redis"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Engine struct {
	App                simple_router.Service
	Log                log.Service
	SQSService         sqs.Service
	SNSService         sns.Service
	DynamoDBService    dynamo.Service
	RedisService       redis.Service
	RepositoriesConfig map[string]interface{}
	UsesCasesConfig    map[string]interface{}
	HandlerConfig      map[string]interface{}
}
