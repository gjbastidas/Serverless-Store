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

func Send(statusCode int, data string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		StatusCode: statusCode,
		Body:       data,
	}, nil
}

func SendOK(aR *APIResponse) (events.APIGatewayProxyResponse, error) {
	log.Info().Msg(aR.LogMessage)
	return Send(aR.StatusCode, aR.Data)
}

func SendErr(aR *APIResponse) (events.APIGatewayProxyResponse, error) {
	log.Error().Msg(aR.LogMessage)
	return Send(aR.StatusCode, aR.Data)
}
