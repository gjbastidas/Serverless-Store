package config

type Cfg struct {
	AWSRegion     string `envconfig:"AWS_REGION" default:"us-east-2"`
	DateString    string `envconfig:"DATESTRING" default:"07-20-2018"`
	ProductsTable string `envconfig:"PRODUCTS_TABLE"`
}
