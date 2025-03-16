package list_items

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	cl "github.com/uala-challenge/simple-toolkit/pkg/client/dynamo/mock"
	log "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestApplyBatchGetItemError(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	// Simulamos error en BatchGetItem
	cli.On("BatchGetItem", mock.Anything, mock.Anything).
		Return(nil, errors.New("error al ejecutar BatchGetItem"))

	l.On("WrapError", mock.Anything, "error al ejecutar BatchGetItem").
		Return(errors.New("error al ejecutar BatchGetItem"))

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	items := []map[string]types.AttributeValue{
		{"PK": &types.AttributeValueMemberS{Value: "123"}},
	}

	rsp, err := cliente.Apply(context.TODO(), items, "table")

	assert.Error(t, err)
	assert.Nil(t, rsp)
	cli.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestApplyUnprocessedKeysRetry(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("BatchGetItem", mock.Anything, mock.Anything).
		Return(&dynamodb.BatchGetItemOutput{
			Responses:       map[string][]map[string]types.AttributeValue{},
			UnprocessedKeys: map[string]types.KeysAndAttributes{"table": {Keys: []map[string]types.AttributeValue{{"PK": &types.AttributeValueMemberS{Value: "456"}}}}},
		}, nil).Once()

	cli.On("BatchGetItem", mock.Anything, mock.Anything).
		Return(&dynamodb.BatchGetItemOutput{
			Responses: map[string][]map[string]types.AttributeValue{
				"table": {{"PK": &types.AttributeValueMemberS{Value: "456"}}},
			},
			UnprocessedKeys: map[string]types.KeysAndAttributes{},
		}, nil).Once()

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	items := []map[string]types.AttributeValue{
		{"PK": &types.AttributeValueMemberS{Value: "123"}},
	}

	start := time.Now()
	rsp, err := cliente.Apply(context.TODO(), items, "table")
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.Len(t, rsp, 1)

	cli.AssertExpectations(t)

	assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(100))
}

func TestApplySuccess(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("BatchGetItem", mock.Anything, mock.Anything).
		Return(&dynamodb.BatchGetItemOutput{
			Responses: map[string][]map[string]types.AttributeValue{
				"table": {
					{"PK": &types.AttributeValueMemberS{Value: "123"}},
					{"PK": &types.AttributeValueMemberS{Value: "456"}},
				},
			},
			UnprocessedKeys: map[string]types.KeysAndAttributes{},
		}, nil).Once()

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	items := []map[string]types.AttributeValue{
		{"PK": &types.AttributeValueMemberS{Value: "123"}},
	}

	rsp, err := cliente.Apply(context.TODO(), items, "table")

	assert.NoError(t, err)
	assert.NotNil(t, rsp)
	assert.Len(t, rsp, 2)

	cli.AssertExpectations(t)
}
