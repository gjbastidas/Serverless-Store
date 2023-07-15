package utils

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

type APIResponse struct {
	StatusCode int
	Message    string
	Data       string
}

func Send(statusCode int, data string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       data,
	}
}

func SendOK(aR *APIResponse) events.APIGatewayProxyResponse {
	log.Info().Msg(aR.Message)
	return Send(aR.StatusCode, aR.Data)
}

func SendErr(aR *APIResponse) events.APIGatewayProxyResponse {
	log.Error().Msg(aR.Message)
	return Send(aR.StatusCode, aR.Data)
}
