package config

type Cfg struct {
	AWSRegion     string `envconfig:"AWS_REGION" default:"us-east-2"`
	ProductsTable string `envconfig:"PRODUCTS_TABLE" default:"products"`
}
