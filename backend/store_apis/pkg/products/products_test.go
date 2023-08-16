package products

import (
	"context"
	"errors"
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
	body := `
	{
		"name": "valid product",
  	"description": "valid product description"
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

	resp, err := p.createOneProduct(context.TODO(), req, cfg, awsSvc)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Contains(t, resp.Body, "successfully created product with id:")
}

func Test_CreateOneProduct_ReturnError(t *testing.T) {
	subtests := []struct {
		name          string
		body          string
		expected      int
		expectedError string
	}{
		{
			name: "invalid_product_name",
			body: `
				{
					"name": "",
					"description": "invalid product description"
				}
			`,
			expected:      http.StatusBadRequest,
			expectedError: "error product validation",
		},
		{
			name: "invalid_product_description",
			body: `
				{
					"name": "invalid product",
					"description": ""
				}
			`,
			expected:      http.StatusBadRequest,
			expectedError: "error product validation",
		},
		{
			name: "error_putting_item",
			body: `
				{
					"name": "valid product name",
					"description": "valid product description"
				}
			`,
			expected:      http.StatusInternalServerError,
			expectedError: "error putting item",
		},
	}

	cfg := new(config.Cfg)
	os.Setenv("PRODUCTS_TABLE", "test")
	err := envconfig.Process("", cfg)
	assert.NoError(t, err)

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			req := events.APIGatewayProxyRequest{
				Headers:    map[string]string{"content-type": "application/json"},
				Resource:   "/products",
				Path:       "/products",
				HTTPMethod: http.MethodPost,
				Body:       st.body,
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDdbClient := mock_aws_services.NewMockDynamoDBClientAPI(ctrl)

			switch st.expectedError {
			case "error putting item":
				mockDdbClient.
					EXPECT().
					PutItem(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("error putting item: some detailed error"))
			default:
				mockDdbClient.
					EXPECT().
					PutItem(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					AnyTimes() // expects zero or more calls
			}

			awsSvc := &aws_services.AWS{
				DDBClient: mockDdbClient,
			}

			p := new(Product)
			resp, err := p.createOneProduct(context.TODO(), req, cfg, awsSvc)
			assert.NoError(t, err)

			assert.Equal(t, st.expected, resp.StatusCode)
			assert.Contains(t, resp.Body, st.expectedError)
		})
	}
}
