package get_item

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	cl "github.com/uala-challenge/simple-toolkit/pkg/client/dynamo/mock"
	log "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestGetItemNotFound(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
		Item:           nil,
		ResultMetadata: middleware.Metadata{},
	}, nil)

	l.On("WrapError", nil, "item no encontrado").Return(errors.New("item no encontrado"))

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})
	rps, err := cliente.Apply(context.TODO(), map[string]interface{}{}, "table")
	assert.Error(t, err)
	assert.Nil(t, rps)

}

func TestGetItemDynamoError(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("GetItem", mock.Anything, mock.Anything).Return(nil, errors.New("error de conexión con DynamoDB"))

	l.On("WrapError", mock.Anything, "error al obtener el item").Return(errors.New("error al obtener el item"))

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})
	rps, err := cliente.Apply(context.TODO(), map[string]interface{}{}, "table")

	assert.Error(t, err)
	assert.Nil(t, rps)
}

func TestGetItemSuccess(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	expectedItem := map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{Value: "123"},
		"SK": &types.AttributeValueMemberS{Value: "abc"},
	}

	cli.On("GetItem", mock.Anything, mock.Anything).Return(&dynamodb.GetItemOutput{
		Item:           expectedItem,
		ResultMetadata: middleware.Metadata{},
	}, nil)

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})
	rps, err := cliente.Apply(context.TODO(), map[string]interface{}{}, "table")

	assert.NoError(t, err)
	assert.NotNil(t, rps)
	assert.Equal(t, expectedItem, rps)
}
