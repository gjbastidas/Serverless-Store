package handlers

import (
	"context"
	"errors"
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
		return utils.SendError(500, err), err
	}

	// set clients
	awsSvc, err := aws_services.NewAWS(cfg.AWSRegion)
	if err != nil {
		return utils.SendError(500, err), err
	}

	switch request.HTTPMethod {
	case http.MethodPost:
		return products.CreateProduct(ctx, request, cfg, awsSvc)
	case http.MethodGet:
		return products.ReadProduct(ctx, request)
	case http.MethodPut:
		return products.UpdateProduct(ctx, request)
	case http.MethodDelete:
		return products.DeleteProduct(ctx, request)
	default:
		err := errors.New("method not defined")
		return utils.SendError(500, err), err
	}
}

// TODO: make handlers for orders and baskets
