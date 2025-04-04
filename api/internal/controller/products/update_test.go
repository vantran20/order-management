package products

import (
	"context"
	"errors"
	"testing"

	"omg/api/internal/model"
	"omg/api/internal/repository"
	"omg/api/internal/repository/inventory"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestImpl_Update(t *testing.T) {
	type args struct {
		input            model.UpdateProductInput
		existingProduct  model.Product
		getProductErr    error
		updateProductOut model.Product
		updateProductErr error
		expectedResult   model.Product
		expectedErr      error
	}

	tcs := map[string]args{
		"success": {
			input: model.UpdateProductInput{
				ID:    123,
				Name:  "Updated Name",
				Desc:  "Updated Description",
				Price: 100.0,
				Stock: 50,
			},
			existingProduct: model.Product{
				ID:     123,
				Name:   "Old Name",
				Status: "active",
			},
			updateProductOut: model.Product{
				ID:          123,
				Name:        "Updated Name",
				Description: "Updated Description",
				Price:       100.0,
				Stock:       50,
				Status:      "active",
			},
			expectedResult: model.Product{
				ID:          123,
				Name:        "Updated Name",
				Description: "Updated Description",
				Price:       100.0,
				Stock:       50,
				Status:      "active",
			},
		},
		"product not found": {
			input: model.UpdateProductInput{
				ID: 123,
			},
			getProductErr: inventory.ErrProductNotFound,
			expectedErr:   ErrNotFound,
		},
		"unexpected get error": {
			input: model.UpdateProductInput{
				ID: 123,
			},
			getProductErr: errors.New("db error"),
			expectedErr:   errors.New("db error"),
		},
		"unexpected update error": {
			input: model.UpdateProductInput{
				ID:    123,
				Name:  "Name",
				Price: 10,
			},
			existingProduct: model.Product{
				ID:     123,
				Status: "active",
			},
			updateProductErr: errors.New("update error"),
			expectedErr:      errors.New("update error"),
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mockInv := &inventory.MockRepository{}
			mockRepo := &repository.MockRegistry{}

			// Setup expected calls
			mockInv.On("GetProductByID", mock.Anything, tc.input.ID).
				Return(tc.existingProduct, tc.getProductErr)

			if tc.getProductErr == nil {
				mockInv.On("UpdateProduct", mock.Anything, mock.MatchedBy(func(p model.Product) bool {
					return p.ID == tc.input.ID
				})).Return(tc.updateProductOut, tc.updateProductErr)
			}

			mockRepo.On("Inventory").Return(mockInv)

			svc := impl{repo: mockRepo}

			// Execute
			result, err := svc.Update(context.Background(), tc.input)

			// Assert
			if tc.expectedErr != nil {
				require.EqualError(t, err, tc.expectedErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedResult, result)
			}

			mockInv.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
