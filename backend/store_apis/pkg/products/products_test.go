package products

import (
	"fmt"
	"net/http"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCreateProduct(t *testing.T) {
	subtests := []struct {
		name             string
		expectedResponse string
		expectedCode     int64
	}{
		{
			name:             "happy_path",
			expectedCode:     http.StatusCreated,
			expectedResponse: "expected",
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			fmt.Print("not implemented")
		})
	}
}

func TestReadProduct(t *testing.T) {
	t.Run("read product", func(t *testing.T) {
		fmt.Print("not implemented")
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("update product", func(t *testing.T) {
		fmt.Print("not implemented")
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("delete product", func(t *testing.T) {
		fmt.Print("not implemented")
	})
}
