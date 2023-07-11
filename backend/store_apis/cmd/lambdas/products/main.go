package main

import (
	"store_apis/pkg/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handlers.ProductsHandler)
}
