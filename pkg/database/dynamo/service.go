package dynamo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/uala-challenge/simple-toolkit/pkg/utilities/log"
)

type service struct {
	client *dynamodb.Client
	config Config
	logger log.Service
}

var _ Service = (*service)(nil)

func NewService(acf aws.Config, cfg Config, logger log.Service) *service {
	options := dynamodb.Options{
		Region: acf.Region,
	}

	if cfg.Endpoint != "" {
		options.BaseEndpoint = aws.String(cfg.Endpoint)
	}

	client := dynamodb.New(options)

	return &service{
		client: client,
		config: cfg,
		logger: logger,
	}
}

func (s *service) PutItem(ctx context.Context, item map[string]interface{}) error {
	if s == nil {
		return errors.New(ErrorServiceNotEnabled)
	}

	av, err := marshalMap(item)
	if err != nil {
		return fmt.Errorf("error al serializar el item: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &s.config.TableName,
		Item:      av,
	})
	if err != nil {
		s.logger.Error(ctx, err, "Error al insertar item en DynamoDB", nil)
		return err
	}

	s.logger.Info(ctx, "Ítem insertado con éxito en DynamoDB", map[string]interface{}{
		"table": s.config.TableName,
	})
	return nil
}

func (s *service) GetItem(ctx context.Context, key map[string]interface{}) (map[string]interface{}, error) {
	if s == nil {
		return nil, errors.New(ErrorServiceNotEnabled)
	}

	av, err := marshalMap(key)
	if err != nil {
		return nil, fmt.Errorf(SerializationError, err)
	}

	result, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &s.config.TableName,
		Key:       av,
	})
	if err != nil {
		s.logger.Error(ctx, err, "Error al obtener item de DynamoDB", nil)
		return nil, err
	}

	if result.Item == nil {
		s.logger.Warn(ctx, "El resultado es nulo", map[string]interface{}{
			"table": s.config.TableName,
		})
		return nil, nil
	}

	return unmarshalMap(result.Item), nil
}

func (s *service) UpdateItem(ctx context.Context, key map[string]interface{}, update map[string]interface{}) error {
	if s == nil {
		return errors.New(ErrorServiceNotEnabled)
	}

	keyAV, err := marshalMap(key)
	if err != nil {
		return fmt.Errorf(SerializationError, err)
	}

	updateAV, err := marshalMap(update)
	if err != nil {
		return fmt.Errorf("error al serializar los datos de actualización: %w", err)
	}

	_, err = s.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: &s.config.TableName,
		Key:       keyAV,
		AttributeUpdates: map[string]types.AttributeValueUpdate{
			"data": {
				Value:  updateAV["data"],
				Action: types.AttributeActionPut,
			},
		},
	})
	if err != nil {
		s.logger.Error(ctx, err, "Error al actualizar item en DynamoDB", nil)
		return err
	}

	s.logger.Info(ctx, "Ítem actualizado con éxito en DynamoDB", nil)
	return nil
}

func (s *service) DeleteItem(ctx context.Context, key map[string]interface{}) error {
	if s == nil {
		return errors.New(ErrorServiceNotEnabled)
	}

	av, err := marshalMap(key)
	if err != nil {
		return fmt.Errorf(SerializationError, err)
	}

	_, err = s.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &s.config.TableName,
		Key:       av,
	})
	if err != nil {
		s.logger.Error(ctx, err, "Error al eliminar item de DynamoDB", nil)
		return err
	}

	s.logger.Info(ctx, "Ítem eliminado con éxito en DynamoDB", nil)
	return nil
}

func marshalMap(input map[string]interface{}) (map[string]types.AttributeValue, error) {
	jsonData, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var av map[string]types.AttributeValue
	err = json.Unmarshal(jsonData, &av)
	return av, err
}

func unmarshalMap(input map[string]types.AttributeValue) map[string]interface{} {
	jsonData, _ := json.Marshal(input)
	var output map[string]interface{}
	_ = json.Unmarshal(jsonData, &output)
	return output
}
