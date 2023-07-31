package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Item struct {
	Id           string `dynamodbav:"id"`
	DateModified int64  `dynamodbav:"dateModified"`
	Name         string `dynamodbav:"name"`
	Description  string `dynamodbav:"description"`
}

func PostHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return createProduct(ctx, request, cfg, awsSvc)
}

func GetHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return readProduct(ctx, request, cfg, awsSvc)
}

func UpdateHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return updateProduct(ctx, request, cfg, awsSvc)
}

func DeleteHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return deleteProduct(ctx, request, cfg, awsSvc)
}

func createProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	product := new(Product)
	if err := json.NewDecoder(strings.NewReader(request.Body)).Decode(product); err != nil {
		msj := fmt.Sprintf("error decoding request body: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusBadRequest,
			Data:       msj,
			LogMessage: msj,
		})
	}

	item := &Item{
		Id:           uuid.New().String(),
		DateModified: time.Now().UTC().Unix(),
		Name:         product.Name,
		Description:  product.Description,
	}

	avMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		msj := fmt.Sprintf("error mapping attribute values: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Item:      avMap,
	}
	_, err = awsSvc.DDBClient.PutItem(ctx, input)
	if err != nil {
		msj := fmt.Sprintf("error putting item: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	msj := fmt.Sprintf("successfully created product with id: %s", item.Id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: http.StatusCreated,
		Data:       msj,
		LogMessage: msj,
	})
}

func readProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		msj := fmt.Sprint("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusBadRequest,
			Data:       msj,
			LogMessage: msj,
		})
	}

	expr, err := expression.NewBuilder().WithKeyCondition(
		expression.
			Key("id").Equal(expression.Value(id)),
	).Build()
	if err != nil {
		msj := fmt.Sprintf("error building query expression: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(cfg.ProductsTable),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(1), // expecting one record only
	}

	queryOutput, err := awsSvc.DDBClient.Query(ctx, queryInput)
	if err != nil {
		msj := fmt.Sprintf("error query item: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	if len(queryOutput.Items) == 0 {
		msj := fmt.Sprintf("no entries found with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	item := new(Item)
	err = attributevalue.UnmarshalMap(queryOutput.Items[0], item)
	if err != nil {
		msj := fmt.Sprintf("error unmarshalling query output: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	out, err := json.Marshal(item)
	if err != nil {
		msj := fmt.Sprintf("error marshalling item: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	return utils.SendOK(&utils.APIResponse{
		StatusCode: http.StatusOK,
		Data:       string(out),
		LogMessage: string(out),
	})
}

func updateProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		msj := fmt.Sprint("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusBadRequest,
			Data:       msj,
			LogMessage: msj,
		})
	}

	product := new(Product)
	err := json.NewDecoder(strings.NewReader(request.Body)).Decode(product)
	if err != nil {
		msj := fmt.Sprintf("error decoding request body: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusBadRequest,
			Data:       msj,
			LogMessage: msj,
		})
	}

	expr, err := expression.NewBuilder().WithUpdate(
		expression.
			Set(expression.Name("name"), expression.Value(product.Name)).
			Set(expression.Name("description"), expression.Value(product.Description)),
	).WithCondition(
		expression.
			AttributeExists(expression.Name("id")),
	).Build()
	if err != nil {
		msj := fmt.Sprintf("error building update expression: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	}
	_, err = awsSvc.DDBClient.UpdateItem(ctx, updateInput)
	if err != nil {
		msj := fmt.Sprintf("error updating item with id: %v. error: %v", id, err)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	msj := fmt.Sprintf("product with id: %v, was successfully updated", id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Data:       msj,
		LogMessage: msj,
	})
}

func deleteProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	if len(id) == 0 {
		msj := fmt.Sprintf("empty id on path params")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusBadRequest,
			Data:       msj,
			LogMessage: msj,
		})
	}

	expr, err := expression.NewBuilder().WithCondition(
		expression.
			AttributeExists(expression.Name("id")),
	).Build()
	if err != nil {
		msj := fmt.Sprintf("error building condition expression: %v", err.Error())
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	deleteInput := &dynamodb.DeleteItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ConditionExpression:       expr.Condition(),
	}
	_, err = awsSvc.DDBClient.DeleteItem(ctx, deleteInput)
	if err != nil {
		msj := fmt.Sprintf("error deleting item with id: %v", id)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	msj := fmt.Sprintf("product with id: %v, was successfully deleted", id)
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Data:       msj,
		LogMessage: msj,
	})
}
