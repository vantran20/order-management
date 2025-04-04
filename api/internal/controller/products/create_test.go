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

func TestImpl_CreateProduct(t *testing.T) {
	type arg struct {
		givenInput              model.CreateProductInput
		mockGetProductByNameOut model.Product
		mockGetProductByNameErr error
		mockCreateProductOut    model.Product
		mockCreateProductErr    error
		expRepoMockCalled       bool
		expResult               model.Product
		expErr                  error
	}

	tcs := map[string]arg{
		"success": {
			givenInput: model.CreateProductInput{
				Name:  "New Product",
				Desc:  "Product description",
				Price: 99.99,
				Stock: 100,
			},
			mockGetProductByNameErr: inventory.ErrProductNotFound,
			mockCreateProductOut: model.Product{
				ID:          1,
				Name:        "New Product",
				Description: "Product description",
				Status:      model.ProductStatusActive,
				Price:       99.99,
				Stock:       100,
			},
			mockCreateProductErr: nil,
			expRepoMockCalled:    true,
			expResult: model.Product{
				ID:          1,
				Name:        "New Product",
				Description: "Product description",
				Status:      model.ProductStatusActive,
				Price:       99.99,
				Stock:       100,
			},
			expErr: nil,
		},
		"product_already_exists": {
			givenInput: model.CreateProductInput{
				Name:  "Existing Product",
				Desc:  "Product description",
				Price: 99.99,
				Stock: 100,
			},
			mockGetProductByNameOut: model.Product{
				ID:          1,
				Name:        "Existing Product",
				Description: "Product description",
				Status:      model.ProductStatusActive,
				Price:       99.99,
				Stock:       100,
			},
			mockGetProductByNameErr: nil,
			expRepoMockCalled:       true,
			expResult:               model.Product{},
			expErr:                  ErrProductAlreadyExists,
		},
		"get_product_error": {
			givenInput: model.CreateProductInput{
				Name:  "Test Product",
				Desc:  "Product description",
				Price: 99.99,
				Stock: 100,
			},
			mockGetProductByNameErr: errors.New("database error"),
			expRepoMockCalled:       true,
			expResult:               model.Product{},
			expErr:                  errors.New("database error"),
		},
		"create_product_error": {
			givenInput: model.CreateProductInput{
				Name:  "New Product",
				Desc:  "Product description",
				Price: 99.99,
				Stock: 100,
			},
			mockGetProductByNameErr: inventory.ErrProductNotFound,
			mockCreateProductErr:    errors.New("failed to create product"),
			expRepoMockCalled:       true,
			expResult:               model.Product{},
			expErr:                  errors.New("failed to create product"),
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			inventoryRepo := inventory.MockRepository{}

			if tc.expRepoMockCalled {
				// Mock GetProductByName call
				inventoryRepo.On("GetProductByName", mock.Anything, tc.givenInput.Name).
					Return(tc.mockGetProductByNameOut, tc.mockGetProductByNameErr)

				// If we expect to reach the CreateProduct call (when product doesn't exist)
				if tc.mockGetProductByNameErr == inventory.ErrProductNotFound {
					// Match the product model
					inventoryRepo.On("CreateProduct", mock.Anything, mock.MatchedBy(func(p model.Product) bool {
						return p.Name == tc.givenInput.Name &&
							p.Description == tc.givenInput.Desc &&
							p.Price == tc.givenInput.Price &&
							p.Stock == tc.givenInput.Stock &&
							p.Status == model.ProductStatusActive
					})).Return(tc.mockCreateProductOut, tc.mockCreateProductErr)
				}
			}

			repo := repository.MockRegistry{}
			repo.On("Inventory").Return(&inventoryRepo)

			impl := impl{repo: &repo}

			// When:
			result, err := impl.Create(context.Background(), tc.givenInput)

			// Then:
			if tc.expErr != nil {
				require.Error(t, err)
				require.Equal(t, tc.expErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}

			inventoryRepo.AssertExpectations(t)
			repo.AssertExpectations(t)
		})
	}
}
