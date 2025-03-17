package query

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"

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

// Test cuando la consulta devuelve resultados
func TestQuerySuccess(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	expectedItems := []map[string]types.AttributeValue{
		{
			"PK": &types.AttributeValueMemberS{Value: "tweet:123"},
			"SK": &types.AttributeValueMemberS{Value: "user:abc"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "tweet:456"},
			"SK": &types.AttributeValueMemberS{Value: "user:def"},
		},
	}

	cli.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
		Items:            expectedItems,
		LastEvaluatedKey: nil, // No más datos para paginar
		ResultMetadata:   middleware.Metadata{},
	}, nil)

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	query := &dynamodb.QueryInput{
		TableName: aws.String("UalaChallenge"),
	}

	rps, err := cliente.Apply(context.TODO(), query)

	assert.NoError(t, err)
	assert.NotNil(t, rps)
	assert.Equal(t, expectedItems, rps)
}

func TestQueryNotFound(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
		Items:            []map[string]types.AttributeValue{}, // Lista vacía
		LastEvaluatedKey: nil,
		ResultMetadata:   middleware.Metadata{},
	}, nil)

	// No esperamos un error en este caso
	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	query := &dynamodb.QueryInput{
		TableName: aws.String("UalaChallenge"),
	}

	rps, err := cliente.Apply(context.TODO(), query)

	// No debe devolver error, solo una lista vacía
	assert.NoError(t, err)
	assert.Empty(t, rps) // Debe devolver un slice vacío
}

// Test cuando DynamoDB devuelve un error
func TestQueryDynamoError(t *testing.T) {
	cli := cl.NewService(t)
	l := log.NewService(t)

	cli.On("Query", mock.Anything, mock.Anything).Return(nil, errors.New("error de conexión con DynamoDB"))

	l.On("WrapError", mock.Anything, "error al ejecutar Query").Return(errors.New("error al ejecutar Query"))

	cliente := NewService(Dependencies{
		Client: &dynamo.Dynamo{Cliente: cli},
		Log:    l,
	})

	query := &dynamodb.QueryInput{
		TableName: aws.String("UalaChallenge"),
	}

	rps, err := cliente.Apply(context.TODO(), query)

	assert.Error(t, err)
	assert.Nil(t, rps)
}
