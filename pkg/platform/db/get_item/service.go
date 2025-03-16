package get_item

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func (s service) Apply(ctx context.Context, item map[string]interface{}, table string) (map[string]types.AttributeValue, error) {
	key, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, s.log.WrapError(err, "error serializando clave de búsqueda")
	}

	result, err := s.client.Cliente.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &table,
		Key:       key,
	})
	if err != nil {
		return nil, s.log.WrapError(err, "error al obtener el item")
	}

	if result.Item == nil {
		return nil, s.log.WrapError(nil, "item no encontrado")
	}

	return result.Item, nil
}
