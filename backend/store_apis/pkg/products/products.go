package products

import (
	"context"

	"store_apis/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
)

func CreateProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parse json body
	// write to database
	return utils.Send(201, "product created"), nil
}

// TODO: make rest of CRUD functions and tests
