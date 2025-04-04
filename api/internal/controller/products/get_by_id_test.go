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

func TestImpl_GetByID(t *testing.T) {
	type arg struct {
		givenID        int64
		mockProduct    model.Product
		mockErr        error
		expReposCalled bool
		expErr         error
	}

	tcs := map[string]arg{
		"success": {
			givenID: 123,
			mockProduct: model.Product{
				ID:     123,
				Status: model.ProductStatusActive,
			},
			expReposCalled: true,
		},
		"product_not_found": {
			givenID:        123,
			mockErr:        inventory.ErrProductNotFound,
			expReposCalled: true,
			expErr:         ErrNotFound,
		},
		"database_error": {
			givenID:        123,
			mockErr:        errors.New("database error"),
			expReposCalled: true,
			expErr:         errors.New("database error"),
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			invRepo := &inventory.MockRepository{}
			if tc.expReposCalled {
				invRepo.On("GetProductByID", mock.Anything, tc.givenID).Return(tc.mockProduct, tc.mockErr)
			}

			mockRepo := &repository.MockRegistry{}
			mockRepo.On("Inventory").Return(invRepo)

			impl := New(mockRepo)

			// When:
			product, err := impl.GetByID(context.Background(), tc.givenID)

			// Then:
			if tc.expErr != nil {
				require.Equal(t, tc.expErr, pkgerrors.Cause(err))
				require.Equal(t, model.Product{}, product)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.mockProduct, product)
			}

			invRepo.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
