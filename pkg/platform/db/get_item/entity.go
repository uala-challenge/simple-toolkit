package get_item

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Apply(ctx context.Context, item map[string]interface{}, table string) (map[string]types.AttributeValue, error)
}

type Dependencies struct {
	Client *dynamo.Dynamo
	Log    log.Service
}
