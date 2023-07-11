package utils

import "github.com/aws/aws-lambda-go/events"

func Send(statusCode int, data string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       data,
	}
}
