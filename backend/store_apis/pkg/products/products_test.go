package products

import (
	"context"
	"net/http"
	"os"
	"store_apis/pkg/config"
	"testing"

	aws_services "store_apis/pkg/aws"
	mock_aws_services "store_apis/pkg/aws/mocks"

	"github.com/aws/aws-lambda-go/events"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_CreateOneProduct_ReturnOK(t *testing.T) {
	ctx := context.TODO()

	body := `
	{
		"name": "one product",
  	"description": "one product description"
	}
	`

	req := events.APIGatewayProxyRequest{
		Headers:    map[string]string{"content-type": "application/json"},
		Resource:   "/products",
		Path:       "/products",
		HTTPMethod: http.MethodPost,
		Body:       body,
	}

	cfg := new(config.Cfg)
	os.Setenv("PRODUCTS_TABLE", "test")
	err := envconfig.Process("", cfg)
	assert.NoError(t, err)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDdbClient := mock_aws_services.NewMockDynamoDBClientAPI(ctrl)

	mockDdbClient.
		EXPECT().
		PutItem(gomock.Any(), gomock.Any()).
		Return(nil, nil)

	awsSvc := &aws_services.AWS{
		DDBClient: mockDdbClient,
	}

	p := new(Product)
	resp, err := p.createOneProduct(ctx, req, cfg, awsSvc)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, resp.Body, "successfully created product with id:")
}
