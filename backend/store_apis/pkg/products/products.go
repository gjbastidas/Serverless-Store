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

func PostHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return createProduct(ctx, request, cfg, awsSvc)
}

func GetHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return readProduct(ctx, request, cfg, awsSvc)
}

func UpdateHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return updateProduct(ctx, request, cfg, awsSvc)
}

func DeleteHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return deleteProduct(ctx, request, cfg, awsSvc)
}

func createProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	var product Product
	if err := json.NewDecoder(strings.NewReader(request.Body)).Decode(&product); err != nil {
		newErr := fmt.Errorf("error decoding request body: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	currentTime := time.Now().UTC().Format(time.RFC3339)
	currentDateTime, err := time.Parse(time.RFC3339, currentTime)
	if err != nil {
		newErr := fmt.Errorf("error parsing datetime: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}
	item := Item{
		Id:           uuid.New().String(),
		DateModified: currentDateTime.String(),
		Name:         product.Name,
		Description:  product.Description,
	}

	avMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		newErr := fmt.Errorf("error mapping attribute values: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Item:      avMap,
	}
	_, err = awsSvc.DDBClient.PutItem(ctx, input)
	if err != nil {
		newErr := fmt.Errorf("error putting item: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	msj := fmt.Sprintf("successfully created product with id: %s", item.Id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 201,
		Data:       msj,
		LogMessage: msj,
	}), nil
}

func readProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		err := errors.New("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Data:       "",
			LogMessage: err.Error(),
		}), err
	}

	keyCondition := expression.
		Key("Id").Equal(expression.Value(id))
	keyExpr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		newErr := fmt.Errorf("error building query expression: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(cfg.ProductsTable),
		KeyConditionExpression:    keyExpr.KeyCondition(),
		ExpressionAttributeNames:  keyExpr.Names(),
		ExpressionAttributeValues: keyExpr.Values(),
	}

	queryOutput, err := awsSvc.DDBClient.Query(ctx, queryInput)
	if err != nil {
		newErr := fmt.Errorf("error query item: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), err
	}

	if len(queryOutput.Items) == 0 {
		err := fmt.Errorf("no entries found with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: err.Error(),
		}), err
	} else if len(queryOutput.Items) > 1 {
		err := fmt.Errorf("duplicate entries found for id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: err.Error(),
		}), err
	}

	items := []Item{}
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &items)
	if err != nil {
		newErr := fmt.Errorf("error unmarshalling query output: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), err
	}

	out, err := json.Marshal(items[0])
	if err != nil {
		newErr := fmt.Errorf("error marshalling items: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), err
	}

	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Data:       string(out),
		LogMessage: string(out),
	}), nil
}

func updateProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		err := errors.New("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Data:       "",
			LogMessage: err.Error(),
		}), err
	}

	keyCondition := expression.Key("Id").Equal(expression.Value(id))
	keyExpr, err := expression.NewBuilder().WithKeyCondition(keyCondition).Build()
	if err != nil {
		newErr := fmt.Errorf("error building query expression: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(cfg.ProductsTable),
		KeyConditionExpression:    keyExpr.KeyCondition(),
		ExpressionAttributeNames:  keyExpr.Names(),
		ExpressionAttributeValues: keyExpr.Values(),
	}

	queryOutput, err := awsSvc.DDBClient.Query(ctx, queryInput)
	if err != nil {
		newErr := fmt.Errorf("error query item: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	if len(queryOutput.Items) == 0 {
		err := fmt.Errorf("no entries found with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: err.Error(),
		}), err
	} else if len(queryOutput.Items) > 1 {
		err := fmt.Errorf("duplicate entries found for id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: err.Error(),
		}), err
	}

	items := []Item{}
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &items)
	if err != nil {
		newErr := fmt.Errorf("error unmarshalling query output: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	key := Key{
		Id:           items[0].Id,
		DateModified: items[0].DateModified,
	}
	avMap, err := attributevalue.MarshalMap(key)
	if err != nil {
		newErr := fmt.Errorf("error mapping attribute values: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	var product Product
	err = json.NewDecoder(strings.NewReader(request.Body)).Decode(&product)
	if err != nil {
		newErr := fmt.Errorf("error decoding request body: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	updateCondition := expression.
		Set(expression.Name("Name"), expression.Value(product.Name)).
		Set(expression.Name("Description"), expression.Value(product.Description))
	updateExpr, err := expression.NewBuilder().WithUpdate(updateCondition).Build()
	if err != nil {
		newErr := fmt.Errorf("error building update expression: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(cfg.ProductsTable),
		Key:                       avMap,
		ExpressionAttributeNames:  updateExpr.Names(),
		ExpressionAttributeValues: updateExpr.Values(),
		UpdateExpression:          updateExpr.Update(),
	}
	_, err = awsSvc.DDBClient.UpdateItem(ctx, updateInput)
	if err != nil {
		newErr := fmt.Errorf("error updating item with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	msj := fmt.Sprintf("product with id: %v, was successfully updated", id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Data:       msj,
		LogMessage: msj,
	}), nil
}

func deleteProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		err := errors.New("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 400,
			Data:       "",
			LogMessage: err.Error(),
		}), err
	}

	keyEx := expression.Key("Id").Equal(expression.Value(id))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		newErr := fmt.Errorf("error building query expression: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(cfg.ProductsTable),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	queryOutput, err := awsSvc.DDBClient.Query(ctx, queryInput)
	if err != nil {
		newErr := fmt.Errorf("error query item: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	if len(queryOutput.Items) == 0 {
		err := fmt.Errorf("no entries found with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       err.Error(),
			LogMessage: err.Error(),
		}), err
	} else if len(queryOutput.Items) > 1 {
		err := fmt.Errorf("duplicate entries found for id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       err.Error(),
			LogMessage: err.Error(),
		}), err
	}

	items := []Item{}
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &items)
	if err != nil {
		newErr := fmt.Errorf("error unmarshalling query output: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	key := Key{
		Id:           items[0].Id,
		DateModified: items[0].DateModified,
	}
	avMap, err := attributevalue.MarshalMap(key)
	if err != nil {
		newErr := fmt.Errorf("error mapping attribute values: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Key:       avMap,
	}
	_, err = awsSvc.DDBClient.DeleteItem(ctx, deleteInput)
	if err != nil {
		newErr := fmt.Errorf("error deleting item with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Data:       "",
			LogMessage: newErr.Error(),
		}), newErr
	}

	msj := fmt.Sprintf("product with id: %v, was successfully deleted", id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Data:       msj,
		LogMessage: msj,
	}), nil
}
