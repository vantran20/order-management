package orders

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"omg/api/internal/controller/orders"
	"omg/api/internal/model"
	"omg/api/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_UpdateOrderStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockOrderCtrl struct {
		wantCall bool
		id       int64
		status   model.OrderStatus
		output   model.Order
		err      error
	}
	tests := map[string]struct {
		givenID         string
		requestBody     updateOrderRequest
		mockOrderCtrl   mockOrderCtrl
		expStatus       int
		expResponse     map[string]interface{}
		shouldBroadcast bool
	}{
		"successful_update_order_status": {
			givenID: "1",
			requestBody: updateOrderRequest{
				Status: model.OrderStatusPaid.String(),
			},
			mockOrderCtrl: mockOrderCtrl{
				wantCall: true,
				id:       1,
				status:   model.OrderStatusPaid,
				output: model.Order{
					ID:        1,
					UserID:    1,
					Status:    model.OrderStatusPaid,
					TotalCost: 100.0,
					OrderItems: []model.OrderItem{
						{
							ID:        1,
							OrderID:   1,
							ProductID: 1,
							Quantity:  2,
							Price:     50.0,
						},
					},
				},
				err: nil,
			},
			expStatus: http.StatusCreated,
			expResponse: map[string]interface{}{
				"id":         "1",
				"user_id":    "1",
				"total_cost": "100",
				"status":     "PAID",
				"items": []interface{}{
					map[string]interface{}{
						"id":         "1",
						"order_id":   "1",
						"product_id": "1",
						"quantity":   "2",
						"price":      "50",
					},
				},
			},
			shouldBroadcast: true,
		},
		"invalid user_id": {
			givenID: "invalid",
			requestBody: updateOrderRequest{
				Status: model.OrderStatusPaid.String(),
			},
			expStatus: http.StatusBadRequest,
			expResponse: map[string]interface{}{
				"error": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
		"missing user_id": {
			givenID:   "0",
			expStatus: http.StatusBadRequest,
			expResponse: map[string]interface{}{
				"error": "user_id required",
			},
		},
		"product not found": {
			givenID: "1",
			requestBody: updateOrderRequest{
				Status: model.OrderStatusPaid.String(),
			},
			mockOrderCtrl: mockOrderCtrl{
				wantCall: true,
				id:       1,
				status:   model.OrderStatusPaid,
				err:      orders.ErrProductNotFound,
			},
			expStatus: http.StatusBadRequest,
			expResponse: map[string]interface{}{
				"error": "product not found",
			},
		},
		"product out of stock": {
			givenID: "1",
			requestBody: updateOrderRequest{
				Status: model.OrderStatusPaid.String(),
			},
			mockOrderCtrl: mockOrderCtrl{
				wantCall: true,
				id:       1,
				status:   model.OrderStatusPaid,
				err:      orders.ErrProductOutOfStock,
			},
			expStatus: http.StatusBadRequest,
			expResponse: map[string]interface{}{
				"error": "product out of stock",
			},
		},
		"internal server error": {
			givenID: "1",
			requestBody: updateOrderRequest{
				Status: model.OrderStatusPaid.String(),
			},
			mockOrderCtrl: mockOrderCtrl{
				wantCall: true,
				id:       1,
				status:   model.OrderStatusPaid,
				err:      errors.New("unexpected error"),
			},
			expStatus: http.StatusInternalServerError,
			expResponse: map[string]interface{}{
				"error": "internal server error",
			},
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			// Create mocks
			mockCtrl := orders.NewMockController(t)
			mockHub := ws.NewMockHub(t)
			handler := NewHandler(mockCtrl, mockHub)

			// Create a test router
			router := gin.New()
			router.PUT("/authenticated/orders/update/:id", handler.UpdateOrderStatus)

			if tc.mockOrderCtrl.wantCall {
				mockCtrl.On("UpdateOrderStatus", mock.Anything, tc.mockOrderCtrl.id, tc.mockOrderCtrl.status).Return(tc.mockOrderCtrl.output, tc.mockOrderCtrl.err)
			}

			if tc.shouldBroadcast {
				// Create a mock hub without starting its goroutine
				mockHub.On("BroadcastMessage", mock.Anything).Return()
			}

			// Create request body
			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/authenticated/orders/update/"+tc.givenID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Verify response
			require.Equal(t, tc.expStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Compare response with expected response
			require.Equal(t, tc.expResponse, response)
		})
	}
}
