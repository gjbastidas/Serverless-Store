package utils

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

func Send(statusCode int, data string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       data,
	}
}

func SendError(statusCode int, err error) events.APIGatewayProxyResponse {
	log.Error().Msg(err.Error())
	return Send(statusCode, err.Error())
}
