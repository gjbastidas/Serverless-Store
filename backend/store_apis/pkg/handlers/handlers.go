package handlers

import (
	"context"
	"errors"
	"net/http"

	"store_apis/pkg/products"
	"store_apis/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
)

func ProductsHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch request.HTTPMethod {
	case http.MethodPost:
		return products.CreateProduct(ctx, request)
	case http.MethodGet:
		return products.ReadProduct(ctx, request)
	case http.MethodPut:
		return products.UpdateProduct(ctx, request)
	case http.MethodDelete:
		return products.DeleteProduct(ctx, request)
	default:
		err := errors.New("method not defined")
		return utils.Send(500, err.Error()), err
	}
}

// TODO: make handlers for orders and baskets
