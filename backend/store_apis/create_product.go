package store_apis

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
)

func createProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// parse json body
	// write to database
	return send(201, "product created"), nil
}
