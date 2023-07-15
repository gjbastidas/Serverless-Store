package products

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"
	"store_apis/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Product struct {
	Id           string    `json:"id,omitempty"`
	Name         string    `json:"name"`
	DateModified time.Time `json:"dateModified,omitempty"`
}

func CreateProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	var product Product

	err := json.NewDecoder(strings.NewReader(request.Body)).Decode(&product)
	if err != nil {
		return utils.SendError(400, err), err
	}

	item, err := attributevalue.MarshalMap(product)
	if err != nil {
		return utils.SendError(500, err), err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(cfg.ProductsTable),
		Item:      item,
	}

	output, err := awsSvc.DDBClient.PutItem(ctx, input)
	if err != nil {
		return utils.SendError(500, err), err
	}

	out, err := json.Marshal(output.Attributes)
	if err != nil {
		return utils.SendError(500, err), err
	}

	return utils.Send(201, string(out)), nil
}

func ReadProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parse json body

	// write to database
	return utils.Send(200, "product returned"), nil
}

func UpdateProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parse json body

	// write to database
	return utils.Send(200, "product updated"), nil
}

func DeleteProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parse json body

	// write to database
	return utils.Send(200, "product deleted"), nil
}
