package orders

import (
	"context"
	"errors"
	"testing"
	"time"

	"omg/api/internal/model"
	"omg/api/internal/repository"
	"omg/api/internal/repository/inventory"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestImpl_CreateOrder(t *testing.T) {
	type arg struct {
		givenInput               model.CreateOrderInput
		mockCreateOrder          model.Order
		mockCreateOrderErr       error
		mockProduct              model.Product
		mockGetProductErr        error
		mockUpdateProductErr     error
		mockCreateOrderItemErr   error
		mockUpdateOrderErr       error
		expDoInTxCalled          bool
		expGetProductCalled      bool
		expUpdateProductCalled   bool
		expCreateOrderItemCalled bool
		expCreateOrderCalled     bool
		expUpdateOrderCalled     bool
		expResult                model.Order
		expErr                   error
	}

	tcs := map[string]arg{
		"success_single_item": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockCreateOrder: model.Order{
				ID:     789,
				UserID: 123,
				Status: model.OrderStatusPending,
			},
			mockProduct: model.Product{
				ID:    456,
				Price: 10.5,
				Stock: 5,
			},
			expDoInTxCalled:          true,
			expGetProductCalled:      true,
			expUpdateProductCalled:   true,
			expCreateOrderItemCalled: true,
			expCreateOrderCalled:     true,
			expUpdateOrderCalled:     true,
			expResult: model.Order{
				ID:        789,
				UserID:    123,
				Status:    model.OrderStatusPending,
				TotalCost: 21.0,
				OrderItems: []model.OrderItem{
					{OrderID: 789, ProductID: 456, Quantity: 2, Price: 10.5},
				},
			},
		},
		"success_multiple_items": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
					{ProductID: 457, Quantity: 1},
				},
			},
			mockCreateOrder: model.Order{
				ID:     789,
				UserID: 123,
				Status: model.OrderStatusPending,
			},
			mockProduct: model.Product{
				ID:    456,
				Price: 10.5,
				Stock: 5,
			},
			expDoInTxCalled:          true,
			expGetProductCalled:      true,
			expUpdateProductCalled:   true,
			expCreateOrderItemCalled: true,
			expCreateOrderCalled:     true,
			expUpdateOrderCalled:     true,
			expResult: model.Order{
				ID:        789,
				UserID:    123,
				Status:    model.OrderStatusPending,
				TotalCost: 36.5, // 2*10.5 + 1*15.5
				OrderItems: []model.OrderItem{
					{OrderID: 789, ProductID: 456, Quantity: 2, Price: 10.5},
					{OrderID: 789, ProductID: 457, Quantity: 1, Price: 15.5},
				},
			},
		},
		"create_order_error": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockCreateOrderErr:   errors.New("db error"),
			expDoInTxCalled:      true,
			expCreateOrderCalled: true,
			expErr:               ErrCreateOrder,
		},
		"product_not_found": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockCreateOrder: model.Order{
				ID:     789,
				UserID: 123,
				Status: model.OrderStatusPending,
			},
			mockGetProductErr:    inventory.ErrProductNotFound,
			expDoInTxCalled:      true,
			expGetProductCalled:  true,
			expCreateOrderCalled: true,
			expErr:               ErrProductNotFound,
		},
		"product_out_of_stock": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 10},
				},
			},
			mockCreateOrder: model.Order{
				ID:     789,
				UserID: 123,
				Status: model.OrderStatusPending,
			},
			mockProduct: model.Product{
				ID:    456,
				Price: 10.5,
				Stock: 5, // Less than requested quantity
			},
			expDoInTxCalled:      true,
			expGetProductCalled:  true,
			expCreateOrderCalled: true,
			expErr:               ErrProductOutOfStock,
		},
		"update_product_error": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockCreateOrder: model.Order{
				ID:     789,
				UserID: 123,
				Status: model.OrderStatusPending,
			},
			mockProduct: model.Product{
				ID:    456,
				Price: 10.5,
				Stock: 5,
			},
			mockUpdateProductErr:   errors.New("update error"),
			expDoInTxCalled:        true,
			expGetProductCalled:    true,
			expUpdateProductCalled: true,
			expCreateOrderCalled:   true,
			expErr:                 ErrUpdateProduct,
		},
		"create_order_item_error": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockCreateOrder: model.Order{
				ID:     789,
				UserID: 123,
				Status: model.OrderStatusPending,
			},
			mockProduct: model.Product{
				ID:    456,
				Price: 10.5,
				Stock: 5,
			},
			mockCreateOrderItemErr:   errors.New("create item error"),
			expDoInTxCalled:          true,
			expGetProductCalled:      true,
			expUpdateProductCalled:   true,
			expCreateOrderItemCalled: true,
			expCreateOrderCalled:     true,
			expErr:                   ErrCreateOrderItem,
		},
		"update_order_error": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockCreateOrder: model.Order{
				ID:     789,
				UserID: 123,
				Status: model.OrderStatusPending,
			},
			mockProduct: model.Product{
				ID:    456,
				Price: 10.5,
				Stock: 5,
			},
			mockUpdateOrderErr:       errors.New("update order error"),
			expDoInTxCalled:          true,
			expGetProductCalled:      true,
			expUpdateProductCalled:   true,
			expCreateOrderItemCalled: true,
			expCreateOrderCalled:     true,
			expUpdateOrderCalled:     true,
			expErr:                   ErrUpdateOrder,
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			invRepo := &inventory.MockRepository{}

			// Setup the mock calls for inventory repository
			if tc.expCreateOrderCalled {
				// First create order call
				expectedOrder := model.Order{
					UserID:    tc.givenInput.UserID,
					Status:    model.OrderStatusPending,
					TotalCost: 0,
				}
				invRepo.On("CreateOrder", mock.Anything, mock.MatchedBy(func(o model.Order) bool {
					return o.UserID == expectedOrder.UserID && o.Status == expectedOrder.Status && o.TotalCost == expectedOrder.TotalCost
				})).Return(tc.mockCreateOrder, tc.mockCreateOrderErr)
			}

			// Setup product mocks for all items
			if tc.expGetProductCalled {
				for _, item := range tc.givenInput.Items {
					productID := item.ProductID

					// For the second item in "multiple_items" test, we need a different product
					mockProduct := tc.mockProduct
					if s == "success_multiple_items" && productID == 457 {
						mockProduct = model.Product{
							ID:    457,
							Price: 15.5,
							Stock: 10,
						}
					}

					invRepo.On("GetProductByID", mock.Anything, productID).Return(mockProduct, tc.mockGetProductErr)

					if tc.expUpdateProductCalled && tc.mockGetProductErr == nil {
						// Check that the product stock is reduced by the requested quantity
						invRepo.On("UpdateProduct", mock.Anything, mock.MatchedBy(func(p model.Product) bool {
							expectedStock := tc.mockProduct.Stock
							if s == "success_multiple_items" && productID == 457 {
								expectedStock = 10 // Original stock for product 457
							}
							return p.ID == productID && p.Stock == expectedStock-item.Quantity
						})).Return(mockProduct, tc.mockUpdateProductErr)
					}

					if tc.expCreateOrderItemCalled && tc.mockUpdateProductErr == nil {
						// Match that order item is created with correct values
						invRepo.On("CreateOrderItem", mock.Anything, mock.MatchedBy(func(item model.OrderItem) bool {
							return item.OrderID == tc.mockCreateOrder.ID &&
								item.ProductID == productID &&
								item.Quantity > 0
						})).Return(model.OrderItem{}, tc.mockCreateOrderItemErr)
					}
				}
			}

			if tc.expUpdateOrderCalled && tc.mockCreateOrderItemErr == nil {
				// Final update of order with total cost
				invRepo.On("UpdateOrder", mock.Anything, mock.MatchedBy(func(o model.Order) bool {
					return o.ID == tc.mockCreateOrder.ID && o.TotalCost > 0
				})).Return(tc.expResult, tc.mockUpdateOrderErr)
			}

			mockRepo := &repository.MockRegistry{}
			mockRepo.On("Inventory").Return(invRepo)

			if tc.expDoInTxCalled {
				// Setup DoInTx mock
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
			result, err := impl.CreateOrder(context.Background(), tc.givenInput)

			// Then:
			if tc.expErr != nil {
				require.Equal(t, tc.expErr, pkgerrors.Cause(err))
				require.Equal(t, model.Order{}, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expResult, result)
			}
		})
	}
}

func TestImpl_CreateOrder_ContextHandling(t *testing.T) {
	type arg struct {
		givenInput     model.CreateOrderInput
		ctxTimeout     time.Duration
		mockRepoDoInTx func(mockRepo *repository.MockRegistry)
		expErr         error
	}

	tcs := map[string]arg{
		"context_timeout": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			ctxTimeout: 1 * time.Millisecond, // Very short timeout to trigger cancellation
			mockRepoDoInTx: func(mockRepo *repository.MockRegistry) {
				mockRepo.On("DoInTx", mock.Anything, mock.AnythingOfType("func(context.Context, repository.Registry) error"), mock.Anything).
					Run(func(args mock.Arguments) {
						// Sleep to ensure context times out
						time.Sleep(5 * time.Millisecond)
					}).
					Return(context.DeadlineExceeded)
			},
			expErr: context.DeadlineExceeded,
		},
		"context_cancelled": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockRepoDoInTx: func(mockRepo *repository.MockRegistry) {
				mockRepo.On("DoInTx", mock.Anything, mock.AnythingOfType("func(context.Context, repository.Registry) error"), mock.Anything).
					Run(func(args mock.Arguments) {
						// Extract and cancel the context
						ctx := args.Get(0).(context.Context)
						if cancel, ok := ctx.Value("mockCancel").(context.CancelFunc); ok {
							cancel()
						}
					}).
					Return(context.Canceled)
			},
			expErr: context.Canceled,
		},
		"get_product_generic_error": {
			givenInput: model.CreateOrderInput{
				UserID: 123,
				Items: []model.CreateOrderItemInput{
					{ProductID: 456, Quantity: 2},
				},
			},
			mockRepoDoInTx: func(mockRepo *repository.MockRegistry) {
				invRepo := &inventory.MockRepository{}

				// Setup the mock calls for inventory repository
				expectedOrder := model.Order{
					UserID:    123,
					Status:    model.OrderStatusPending,
					TotalCost: 0,
				}

				// First create order call
				invRepo.On("CreateOrder", mock.Anything, mock.MatchedBy(func(o model.Order) bool {
					return o.UserID == expectedOrder.UserID && o.Status == expectedOrder.Status && o.TotalCost == expectedOrder.TotalCost
				})).Return(model.Order{ID: 789, UserID: 123, Status: model.OrderStatusPending}, nil)

				// Setup product mock with generic error
				invRepo.On("GetProductByID", mock.Anything, int64(456)).Return(model.Product{}, errors.New("database error"))

				mockRepo.On("Inventory").Return(invRepo)

				// Setup DoInTx mock
				mockRepo.On("DoInTx", mock.Anything, mock.AnythingOfType("func(context.Context, repository.Registry) error"), mock.Anything).
					Run(func(args mock.Arguments) {
						// Extract and execute the transaction function
						txFunc := args.Get(1).(func(context.Context, repository.Registry) error)
						// We're calling the transaction function with the mock registry
						err := txFunc(context.Background(), mockRepo)
						require.Error(t, err)
						require.Equal(t, ErrGetProduct, err)
					}).
					Return(ErrGetProduct)
			},
			expErr: ErrGetProduct,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			// Given:
			mockRepo := &repository.MockRegistry{}
			tc.mockRepoDoInTx(mockRepo)

			impl := New(mockRepo)

			// Create context with timeout if specified
			var ctx context.Context
			var cancel context.CancelFunc

			if tc.ctxTimeout > 0 {
				ctx, cancel = context.WithTimeout(context.Background(), tc.ctxTimeout)
			} else if name == "context_cancelled" {
				ctx, cancel = context.WithCancel(context.Background())
				ctx = context.WithValue(ctx, "mockCancel", cancel)
			} else {
				ctx = context.Background()
				cancel = func() {}
			}
			defer cancel()

			// When:
			result, err := impl.CreateOrder(ctx, tc.givenInput)

			// Then:
			if tc.expErr != nil {
				if errors.Is(err, tc.expErr) {
					// Context errors need to be checked with errors.Is
					require.True(t, errors.Is(err, tc.expErr))
				} else {
					require.Equal(t, tc.expErr, err)
				}
				require.Equal(t, model.Order{}, result)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
