package list_items

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Apply(ctx context.Context, items []map[string]types.AttributeValue, table string) ([]map[string]types.AttributeValue, error)
}

type Dependencies struct {
	Client *dynamo.Dynamo
	Log    log.Service
}
