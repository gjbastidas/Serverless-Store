package aws_services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type AWS struct {
	config    aws.Config
	DDBClient DynamoDBClientAPI
}

func NewAWS(region string) (*AWS, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	// Set respective AWS services clients below
	dynamoDbClient := dynamodb.NewFromConfig(cfg)

	return &AWS{
		config:    cfg,
		DDBClient: dynamoDbClient,
	}, nil
}
