package utils

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog/log"
)

type APIResponse struct {
	StatusCode int
	Data       string
	LogMessage string
}

func Send(statusCode int, data string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       data,
	}
}

func SendOK(aR *APIResponse) events.APIGatewayProxyResponse {
	log.Info().Msg(aR.LogMessage)
	return Send(aR.StatusCode, aR.Data)
}

func SendErr(aR *APIResponse) events.APIGatewayProxyResponse {
	log.Error().Msg(aR.LogMessage)
	return Send(aR.StatusCode, aR.Data)
}
