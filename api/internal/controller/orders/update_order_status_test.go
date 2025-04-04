package orders

import (
	"context"
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

			impl := New(mockRepo)

			// When:
			rs, err := impl.UpdateOrderStatus(context.Background(), tc.givenID, tc.givenStatus)

			// Then:
			if tc.expErr != nil {
				require.ErrorIs(t, err, tc.expErr)
				require.Equal(t, model.Order{}, rs)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.givenID, rs.ID)
				require.Equal(t, tc.givenStatus, rs.Status)
			}

			invRepo.AssertExpectations(t)
			mockRepo.AssertExpectations(t)
		})
	}
}
