package products

import (
	"context"
	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"

	"github.com/aws/aws-lambda-go/events"
)

func PostHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	p := new(Product)
	return p.createProduct(ctx, request, cfg, awsSvc)
}

func GetHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	p := new(Product)
	return p.readProduct(ctx, request, cfg, awsSvc)
}

func UpdateHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	p := new(Product)
	return p.updateProduct(ctx, request, cfg, awsSvc)
}

func DeleteHandler(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error) {
	p := new(Product)
	return p.deleteProduct(ctx, request, cfg, awsSvc)
}
