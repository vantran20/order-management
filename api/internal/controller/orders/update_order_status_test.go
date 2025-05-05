package orders

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

func TestImpl_UpdateOrderStatus(t *testing.T) {
	type arg struct {
		givenID       int64
		givenStatus   model.OrderStatus
		mockOrder     model.Order
		mockGetErr    error
		mockUpdateErr error

		expGetCalled    bool
		expUpdateCalled bool
		expErr          error
	}

	tcs := map[string]arg{
		"success": {
			givenID:     1,
			givenStatus: model.OrderStatusShipped,
			mockOrder: model.Order{
				ID:     1,
				Status: model.OrderStatusPending,
			},
			expGetCalled:    true,
			expUpdateCalled: true,
		},
		"order_not_found_on_get": {
			givenID:      2,
			givenStatus:  model.OrderStatusCancelled,
			mockGetErr:   inventory.ErrOrderNotFound,
			expGetCalled: true,
			expErr:       ErrOrderNotFound,
		},
		"order_not_found_on_update": {
			givenID:     4,
			givenStatus: model.OrderStatusCancelled,
			mockOrder: model.Order{
				ID:     4,
				Status: model.OrderStatusPending,
			},
			mockUpdateErr:   inventory.ErrOrderNotFound,
			expGetCalled:    true,
			expUpdateCalled: true,
			expErr:          ErrOrderNotFound,
		},
		"generic_error_on_get": {
			givenID:      3,
			givenStatus:  model.OrderStatusPaid,
			mockGetErr:   errors.New("database connection error"),
			expGetCalled: true,
			expErr:       errors.New("database connection error"),
		},
		"generic_error_on_update": {
			givenID:     5,
			givenStatus: model.OrderStatusPaid,
			mockOrder: model.Order{
				ID:     5,
				Status: model.OrderStatusPending,
			},
			mockUpdateErr:   errors.New("database error"),
			expGetCalled:    true,
			expUpdateCalled: true,
			expErr:          errors.New("database error"),
		},
		"zero_id_check": {
			givenID:         0,
			givenStatus:     model.OrderStatusPaid,
			mockOrder:       model.Order{},
			expGetCalled:    true,
			expUpdateCalled: true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			invRepo := &inventory.MockRepository{}
			if tc.expGetCalled {
				invRepo.On("GetOrderByID", mock.Anything, tc.givenID).Return(tc.mockOrder, tc.mockGetErr)
			}
			if tc.expUpdateCalled {
				expectedOrder := tc.mockOrder
				expectedOrder.Status = tc.givenStatus
				invRepo.On("UpdateOrder", mock.Anything, expectedOrder).Return(expectedOrder, tc.mockUpdateErr)
			}

			mockRepo := &repository.MockRegistry{}
			mockRepo.On("Inventory").Return(invRepo)

			i := New(mockRepo)

			// When:
			rs, err := i.UpdateOrderStatus(context.Background(), tc.givenID, tc.givenStatus)

			// Then:
			if tc.expErr != nil {
				require.EqualError(t, err, tc.expErr.Error())
				require.Equal(t, model.Order{}, rs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.givenID, rs.ID)
				require.Equal(t, tc.givenStatus, rs.Status)
			}
		})
	}
}
