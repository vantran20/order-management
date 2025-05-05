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

func Test_impl_List(t *testing.T) {
	type arg struct {
		mockProducts   []model.Product
		mockErr        error
		expReposCalled bool
		expErr         error
	}

	tcs := map[string]arg{
		"success": {
			mockProducts: []model.Product{
				{
					ID:     123,
					Status: model.ProductStatusActive,
				},
			},
			expReposCalled: true,
		},
		"empty": {
			mockProducts:   nil,
			mockErr:        nil,
			expReposCalled: true,
			expErr:         nil,
		},
		"database_error": {
			mockProducts:   nil,
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
				invRepo.On("ListProducts", mock.Anything).Return(tc.mockProducts, tc.mockErr)
			}

			mockRepo := &repository.MockRegistry{}
			mockRepo.On("Inventory").Return(invRepo)

			impl := New(mockRepo)

			// When:
			products, err := impl.List(context.Background())

			// Then:
			if tc.expErr != nil {
				require.Equal(t, tc.expErr, pkgerrors.Cause(err))
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.mockProducts, products)
			}
		})
	}
}
