package dynamo

import (
	"context"
)

const (
	ErrorServiceNotEnabled = "DynamoDB Service no habilitado"
	SerializationError     = "error al serializar la clave: %w"
)

type Config struct {
	TableName string `json:"table_name"`
	Region    string `json:"region"`
	Endpoint  string `json:"endpoint"`
}

type Service interface {
	PutItem(ctx context.Context, item map[string]interface{}) error
	GetItem(ctx context.Context, key map[string]interface{}) (map[string]interface{}, error)
	UpdateItem(ctx context.Context, key map[string]interface{}, update map[string]interface{}) error
	DeleteItem(ctx context.Context, key map[string]interface{}) error
}
