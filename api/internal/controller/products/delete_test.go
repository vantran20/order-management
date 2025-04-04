package products

import (
	"context"
	"errors"
	"testing"

	"omg/api/internal/model"
	"omg/api/internal/repository"
	"omg/api/internal/repository/inventory"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestImpl_Delete(t *testing.T) {
	type arg struct {
		givenID         int64
		mockGetProduct  model.Product
		mockGetErr      error
		mockUpdateErr   error
		expGetCalled    bool
		expUpdateCalled bool
		expDoInTxCalled bool
		expErr          error
	}

	tcs := map[string]arg{
		"success": {
			givenID: 123,
			mockGetProduct: model.Product{
				ID:     123,
				Status: model.ProductStatusActive,
			},
			expGetCalled:    true,
			expUpdateCalled: true,
			expDoInTxCalled: true,
		},
		"product_not_found": {
			givenID:         123,
			mockGetErr:      inventory.ErrProductNotFound,
			expGetCalled:    true,
			expDoInTxCalled: true,
			expErr:          ErrNotFound,
		},
		"product_already_deleted": {
			givenID: 123,
			mockGetProduct: model.Product{
				ID:     123,
				Status: model.ProductStatusDeleted,
			},
			expGetCalled:    true,
			expDoInTxCalled: true,
			expErr:          ErrProductDeleted,
		},
		"update_error": {
			givenID: 123,
			mockGetProduct: model.Product{
				ID:     123,
				Status: model.ProductStatusActive,
			},
			mockUpdateErr:   errors.New("update error"),
			expGetCalled:    true,
			expUpdateCalled: true,
			expDoInTxCalled: true,
			expErr:          errors.New("update error"),
		},
		"update_not_found_error": {
			givenID: 123,
			mockGetProduct: model.Product{
				ID:     123,
				Status: model.ProductStatusActive,
			},
			mockUpdateErr:   inventory.ErrProductNotFound,
			expGetCalled:    true,
			expUpdateCalled: true,
			expDoInTxCalled: true,
			expErr:          ErrNotFound,
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			invRepo := &inventory.MockRepository{}
			if tc.expGetCalled {
				invRepo.On("GetProductByID", mock.Anything, tc.givenID).Return(tc.mockGetProduct, tc.mockGetErr)
			}

			if tc.expUpdateCalled {
				// Use mock.MatchedBy to verify the product status was set to deleted
				invRepo.On("UpdateProduct", mock.Anything, mock.MatchedBy(func(p model.Product) bool {
					return p.ID == tc.givenID && p.Status == model.ProductStatusDeleted
				})).Return(tc.mockGetProduct, tc.mockUpdateErr)
			}

			mockRepo := &repository.MockRegistry{}
			mockRepo.On("Inventory").Return(invRepo)

			if tc.expDoInTxCalled {
				// Use mock.AnythingOfType to match the transaction function
				mockRepo.On("DoInTx", mock.Anything, mock.AnythingOfType("func(context.Context, repository.Registry) error"), mock.Anything).
					Run(func(args mock.Arguments) {
						// Extract and execute the transaction function
						txFunc := args.Get(1).(func(context.Context, repository.Registry) error)
						// We're calling the transaction function with the mock registry
						if err := txFunc(context.Background(), mockRepo); err != nil {
							// Do nothing, we're just simulating the transaction
						}
					}).
					Return(tc.expErr)
			}

			impl := New(mockRepo)

			// When:
			err := impl.Delete(context.Background(), tc.givenID)

			// Then:
			if tc.expErr != nil {
				require.Equal(t, tc.expErr, pkgerrors.Cause(err))
			} else {
				require.NoError(t, err)
			}

			invRepo.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
