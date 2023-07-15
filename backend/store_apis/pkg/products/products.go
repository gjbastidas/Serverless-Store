package products

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"
	"store_apis/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
)

type Product struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Item struct {
	Id           string `json:"id"`
	DateModified string `json:"dateModified"`
	Name         string `json:"name"`
	Description  string `json:"description"`
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

	output, err := awsSvc.DDBClient.PutItem(ctx, input)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error putting item: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	out, err := json.Marshal(output.Attributes)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("error marshalling attributes: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	return utils.SendOK(&utils.APIResponse{
		StatusCode: 201,
		Message:    fmt.Sprintf("successfully created product with id: %s", item.Id),
		Data:       string(out),
	}), nil
}

func ReadProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	msj := "product returned"
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Message:    msj,
		Data:       msj,
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

func DeleteProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	msj := "product deleted"
	return utils.SendOK(&utils.APIResponse{
		StatusCode: 200,
		Message:    msj,
		Data:       msj,
	}), nil
}
