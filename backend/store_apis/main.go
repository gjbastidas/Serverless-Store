package store_apis

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func ProductsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodPost:
		return createProduct(ctx, request)
	default:
		err := errors.New("method not defined")
		return send(500, err.Error()), err
	}
}

func send(statusCode int, data string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       data,
	}
}
