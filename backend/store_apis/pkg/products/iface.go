package products

import (
	"context"
	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"

	"github.com/aws/aws-lambda-go/events"
)

type IProduct interface {
	createProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
	readProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
	updateProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
	deleteProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
}
