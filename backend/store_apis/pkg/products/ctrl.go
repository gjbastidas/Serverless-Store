package products

import (
	"context"
	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"

	"github.com/aws/aws-lambda-go/events"
)

func Post(ctx context.Context, request events.APIGatewayProxyRequest, p IProduct, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return p.createOneProduct(ctx, request, cfg, awsSvc)
}

func Get(ctx context.Context, request events.APIGatewayProxyRequest, p IProduct, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return p.readOneProduct(ctx, request, cfg, awsSvc)
}

func Put(ctx context.Context, request events.APIGatewayProxyRequest, p IProduct, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return p.updateOneProduct(ctx, request, cfg, awsSvc)
}

func Delete(ctx context.Context, request events.APIGatewayProxyRequest, p IProduct, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	return p.deleteOneProduct(ctx, request, cfg, awsSvc)
}
