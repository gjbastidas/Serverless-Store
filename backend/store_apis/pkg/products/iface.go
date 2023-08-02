package products

import (
	"context"
	aws_services "store_apis/pkg/aws"
	"store_apis/pkg/config"

	"github.com/aws/aws-lambda-go/events"
)

type IProduct interface {
	createOneProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
	readOneProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
	updateOneProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
	deleteOneProduct(ctx context.Context, request events.APIGatewayProxyRequest, cfg *config.Cfg, awsSvc *aws_services.AWS) (events.APIGatewayProxyResponse, error)
}
