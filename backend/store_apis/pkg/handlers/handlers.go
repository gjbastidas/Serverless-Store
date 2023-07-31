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
)

func ProductsHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx := context.TODO()

	// set config
	cfg := new(config.Cfg)
	err := envconfig.Process("", cfg)
	if err != nil {
		msj := fmt.Sprintf("bad environment configuration: %v", err)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	// set clients
	awsSvc, err := aws_services.NewAWS(cfg.AWSRegion)
	if err != nil {
		msj := fmt.Sprintf("error setting AWS services: %v", err)
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       msj,
			LogMessage: msj,
		})
	}

	switch request.HTTPMethod {
	case http.MethodPost:
		return products.PostHandler(ctx, request, cfg, awsSvc)
	case http.MethodGet:
		return products.GetHandler(ctx, request, cfg, awsSvc)
	case http.MethodPut:
		return products.UpdateHandler(ctx, request, cfg, awsSvc)
	case http.MethodDelete:
		return products.DeleteHandler(ctx, request, cfg, awsSvc)
	default:
		err := errors.New("method not defined")
		return utils.SendErr(&utils.APIResponse{
			StatusCode: http.StatusInternalServerError,
			Data:       err.Error(),
			LogMessage: err.Error(),
		})
	}
}

// TODO: make handlers for orders and baskets
