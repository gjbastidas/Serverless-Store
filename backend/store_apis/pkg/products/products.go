package products

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"
	"store_apis/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Item struct {
	Id           string
	DateModified string
	Name         string
	Description  string
}

type Key struct {
	Id           string
	DateModified string
}

func CreateProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	var product Product

	err := json.NewDecoder(strings.NewReader(request.Body)).Decode(&product)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Message:    fmt.Sprintf("error decoding request body: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	item := Item{
		Id:           uuid.New().String(),
		DateModified: time.Now().Format(cfg.DateString),
		Name:         product.Name,
		Description:  product.Description,
	}

	avMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error mapping attribute values: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Item:      avMap,
	}

	_, err = awsSvc.DDBClient.PutItem(ctx, input)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error putting item: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	msj := fmt.Sprintf("successfully created product with id: %s", item.Id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 201,
		Message:    msj,
		Data:       msj,
	}), nil
}

func ReadProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		err := errors.New("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Message:    err.Error(),
			Data:       err.Error(),
		}), err
	}

	keyEx := expression.Key("Id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error building query expression: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(cfg.ProductsTable),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	queryOutput, err := awsSvc.DDBClient.Query(ctx, queryInput)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error query item: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	items := []Item{}
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &items)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error unmarshalling query output: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	if len(items) > 1 {
		err := fmt.Errorf("duplicate entries found for id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    err.Error(),
			Data:       err.Error(),
		}), err
	}

	out, err := json.Marshal(items[0])
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error marshalling items: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Message:    string(out),
		Data:       string(out),
	}), nil
}

func UpdateProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	msj := "product updated"
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Message:    msj,
		Data:       msj,
	}), nil
}

func DeleteProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		err := errors.New("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Message:    err.Error(),
			Data:       err.Error(),
		}), err
	}

	keyEx := expression.Key("Id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error building query expression: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(cfg.ProductsTable),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	queryOutput, err := awsSvc.DDBClient.Query(ctx, queryInput)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error query item: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	if len(queryOutput.Items) == 0 {
		err := fmt.Errorf("no entries found with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    err.Error(),
			Data:       err.Error(),
		}), err
	} else if len(queryOutput.Items) > 1 {
		err := fmt.Errorf("duplicate entries found for id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    err.Error(),
			Data:       err.Error(),
		}), err
	}

	items := []Item{}
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &items)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error unmarshalling query output: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	key := Key{
		Id:           items[0].Id,
		DateModified: items[0].DateModified,
	}
	avMap, err := attributevalue.MarshalMap(key)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error mapping attribute values: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Key:       avMap,
	}
	_, err = awsSvc.DDBClient.DeleteItem(ctx, deleteInput)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error deleting item with id: %v", id),
			Data:       err.Error(),
		}), err
	}

	msj := fmt.Sprintf("product with id: %v, was successfully deleted", id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Message:    msj,
		Data:       msj,
	}), nil
}
