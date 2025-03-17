package query

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *dynamo.Dynamo
	log    log.Service
}

var _ Service = (*service)(nil)

func NewService(d Dependencies) *service {
	return &service{
		client: d.Client,
		log:    d.Log,
	}
}

func (s service) Apply(ctx context.Context, query *dynamodb.QueryInput) ([]map[string]types.AttributeValue, error) {
	var results []map[string]types.AttributeValue
	for {
		output, err := s.client.Cliente.Query(ctx, query)
		if err != nil {
			return nil, s.log.WrapError(err, "error al ejecutar Query")
		}
		if len(output.Items) > 0 {
			results = append(results, output.Items...)
		}
		if output.LastEvaluatedKey == nil {
			break
		}
		query.ExclusiveStartKey = output.LastEvaluatedKey
	}

	return results, nil
}
