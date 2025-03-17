package viper

import (
	"github.com/uala-challenge/simple-toolkit/pkg/client/rest"
	"sync"

	"github.com/sirupsen/logrus"

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
	Router       simple_router.Config     `json:"router" yaml:"router"`
	Rest         []map[string]rest.Config `json:"rest" yaml:"rest"`
	Log          log.Config               `json:"log" yaml:"log"`
	Aws          AwsConfig                `json:"aws" yaml:"aws"`
	SQS          *sqs.Config              `json:"sqs" yaml:"sqs"`
	SNS          *sns.Config              `json:"sns" yaml:"sns"`
	Dynamo       *dynamo.Config           `json:"dynamo" yaml:"dynamo"`
	Redis        *redis.Config            `json:"redis" yaml:"redis"`
	Repositories map[string]interface{}   `json:"repositories" yaml:"repositories"`
	Cases        map[string]interface{}   `json:"cases" yaml:"cases"`
	Endpoints    map[string]interface{}   `json:"endpoints" yaml:"endpoints"`
}

type AwsConfig struct {
	Region string `json:"region"`
}

type service struct {
	propertyFiles []string
	path          string
	log           *logrus.Logger
}

var (
	once     sync.Once
	instance *service
)
