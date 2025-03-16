package save_item

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"
	cl "github.com/uala-challenge/simple-toolkit/pkg/client/dynamo/mock"
	log "github.com/uala-challenge/simple-toolkit/pkg/utilities/log/mock"
)

func TestAcceptPutItemError(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("PutItem", mock.Anything, mock.Anything).
		Return(nil, errors.New("Error al guardar el item"))

	l.On("WrapError", mock.Anything, "Error al guardar el item").
		Return(errors.New("Error al guardar el item"))

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	err := cliente.Accept(context.TODO(), map[string]interface{}{"PK": "123"}, "table")
	assert.Error(t, err)
	cli.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestAcceptSuccess(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("PutItem", mock.Anything, mock.Anything).
		Return(&dynamodb.PutItemOutput{}, nil)

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	err := cliente.Accept(context.TODO(), map[string]interface{}{"PK": "123"}, "table")
	assert.NoError(t, err)
	cli.AssertExpectations(t)
}
