package list_items

import (
	"context"
	"time"

	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

func (s service) Apply(ctx context.Context, items []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error) {
	batchInput := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			table: {Keys: items},
		},
	}

	var results []map[string]types.AttributeValue

	for attempts := 0; attempts < 3; attempts++ {
		output, err := s.client.Cliente.BatchGetItem(ctx, batchInput)
		if err != nil {
			return nil, s.log.WrapError(err, "error al ejecutar BatchGetItem")
		}

		if len(output.Responses[table]) > 0 {
			if results == nil {
				results = make([]map[string]types.AttributeValue, 0, len(output.Responses[table]))
			}
			results = append(results, output.Responses[table]...)
		}

		if len(output.UnprocessedKeys) == 0 {
			break
		}

		batchInput.RequestItems = output.UnprocessedKeys
		time.Sleep(time.Millisecond * 100)
	}

	return results, nil
}
