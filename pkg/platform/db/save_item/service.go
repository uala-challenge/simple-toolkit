package save_item

import (
	"context"

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

func (s *service) Accept(ctx context.Context, itm map[string]interface{}, table string) error {
	item, err := attributevalue.MarshalMap(itm)
	if err != nil {
		return s.log.WrapError(err, "Error serializando item")
	}
	_, err = s.client.Cliente.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &table,
		Item:      item,
	})
	if err != nil {
		return s.log.WrapError(err, "Error al guardar el item")
	}
	return nil
}
