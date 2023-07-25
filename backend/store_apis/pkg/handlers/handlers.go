package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"
	"store_apis/pkg/products"
	"store_apis/pkg/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

func ProductsHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := context.TODO()

	// set config
	var cfg config.Cfg
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Error().Msg("bad environment configuration")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("bad environment configuration: %v", err.Error()),
			Data:       err.Error(),
		}), err
	}

	// set clients
	awsSvc, err := aws_services.NewAWS(cfg.AWSRegion)
	if err != nil {
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    err.Error(),
			Data:       err.Error(),
		}), err
	}

	switch request.HTTPMethod {
	case http.MethodPost:
		return products.CreateProduct(ctx, request, cfg, awsSvc)
	case http.MethodGet:
		return products.ReadProduct(ctx, request, cfg, awsSvc)
	case http.MethodPut:
		return products.UpdateProduct(ctx, request, cfg, awsSvc)
	case http.MethodDelete:
		return products.DeleteProduct(ctx, request, cfg, awsSvc)
	default:
		err := errors.New("method not defined")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: 500,
			Message:    err.Error(),
			Data:       err.Error(),
		}), err
	}
}

// TODO: make handlers for orders and baskets
