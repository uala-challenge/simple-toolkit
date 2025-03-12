package viper

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/uala-challenge/simple-toolkit/pkg/simplify/simple_router"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
	"sync"
)

type Service interface {
	Apply() (Config, error)
}

type Config struct {
	Router       simple_router.Config   `json:"router"`
	Log          log.Config             `json:"log"`
	AwsConfig    aws.Config             `json:"awsConfig"`
	Repositories map[string]interface{} `json:"repositories"`
	UsesCases    map[string]interface{} `json:"uses_cases"`
	Endpoints    map[string]interface{} `json:"endpoints"`
}

type service struct {
	propertyFiles []string
	path          string
}

var (
	once     sync.Once
	instance *service
)
