package save_item

import (
	"context"

	"github.com/uala-challenge/simple-toolkit/pkg/client/dynamo"

	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type Service interface {
	Accept(ctx context.Context, itm map[string]interface{}, table string) error
}

type Dependencies struct {
	Client *dynamo.Dynamo
	Log    log.Service
}
