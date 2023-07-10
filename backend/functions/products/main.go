package main

import (
	sapis "../../store_apis"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(sapis.ProductsHandler)
}
