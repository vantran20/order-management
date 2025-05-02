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

func TestHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockOrderCtrl struct {
		input  model.CreateOrderInput
		output model.Order
		err    error
	}
	tests := map[string]struct {
		requestBody     createOrderRequest
		mockOrderCtrl   mockOrderCtrl
		expStatus       int
		expResponse     map[string]interface{}
		shouldBroadcast bool
	}{
		"successful order creation": {
			requestBody: createOrderRequest{
				UserID: "1",
				Items: []struct {
					ProductID string `json:"product_id"`
					Quantity  string `json:"quantity"`
				}{
					{
						ProductID: "1",
						Quantity:  "2",
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				input: model.CreateOrderInput{
					UserID: 1,
					Items: []model.CreateOrderItemInput{
						{
							ProductID: 1,
							Quantity:  2,
						},
					},
				},
				output: model.Order{
					ID:        1,
					UserID:    1,
					Status:    model.OrderStatusPending,
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
				"status":     "PENDING",
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
			requestBody: createOrderRequest{
				UserID: "invalid",
				Items: []struct {
					ProductID string `json:"product_id"`
					Quantity  string `json:"quantity"`
				}{
					{
						ProductID: "1",
						Quantity:  "2",
					},
				},
			},
			expStatus: http.StatusBadRequest,
			expResponse: map[string]interface{}{
				"error": "strconv.ParseInt: parsing \"invalid\": invalid syntax",
			},
		},
		"missing user_id": {
			requestBody: createOrderRequest{
				UserID: "0",
				Items: []struct {
					ProductID string `json:"product_id"`
					Quantity  string `json:"quantity"`
				}{
					{
						ProductID: "1",
						Quantity:  "2",
					},
				},
			},
			expStatus: http.StatusBadRequest,
			expResponse: map[string]interface{}{
				"error": "user_id required",
			},
		},
		"product not found": {
			requestBody: createOrderRequest{
				UserID: "1",
				Items: []struct {
					ProductID string `json:"product_id"`
					Quantity  string `json:"quantity"`
				}{
					{
						ProductID: "1",
						Quantity:  "2",
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				input: model.CreateOrderInput{
					UserID: 1,
					Items: []model.CreateOrderItemInput{
						{
							ProductID: 1,
							Quantity:  2,
						},
					},
				},
				err: orders.ErrProductNotFound,
			},
			expStatus: http.StatusConflict,
			expResponse: map[string]interface{}{
				"error": "product not found",
			},
		},
		"product out of stock": {
			requestBody: createOrderRequest{
				UserID: "1",
				Items: []struct {
					ProductID string `json:"product_id"`
					Quantity  string `json:"quantity"`
				}{
					{
						ProductID: "1",
						Quantity:  "2",
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				input: model.CreateOrderInput{
					UserID: 1,
					Items: []model.CreateOrderItemInput{
						{
							ProductID: 1,
							Quantity:  2,
						},
					},
				},
				err: orders.ErrProductOutOfStock,
			},
			expStatus: http.StatusConflict,
			expResponse: map[string]interface{}{
				"error": "product out of stock",
			},
		},
		"internal server error": {
			requestBody: createOrderRequest{
				UserID: "1",
				Items: []struct {
					ProductID string `json:"product_id"`
					Quantity  string `json:"quantity"`
				}{
					{
						ProductID: "1",
						Quantity:  "2",
					},
				},
			},
			mockOrderCtrl: mockOrderCtrl{
				input: model.CreateOrderInput{
					UserID: 1,
					Items: []model.CreateOrderItemInput{
						{
							ProductID: 1,
							Quantity:  2,
						},
					},
				},
				err: errors.New("unexpected error"),
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
			if tc.mockOrderCtrl.input.UserID != 0 {
				mockCtrl.On("CreateOrder", mock.Anything, tc.mockOrderCtrl.input).Return(tc.mockOrderCtrl.output, tc.mockOrderCtrl.err)
			}

			// Set up handler
			h := Handler{
				controller: mockCtrl,
			}
			if tc.shouldBroadcast {
				// Create a mock hub without starting its goroutine
				mockHub := ws.NewMockHub(t)
				mockHub.On("BroadcastMessage", mock.Anything).Return()
				h.wsHub = mockHub
			}

			// Create request body
			body, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			// Create test request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Execute request
			h.Create(c)

			// Verify response
			require.Equal(t, tc.expStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			// Compare response with expected response
			require.Equal(t, tc.expResponse, response)

			// Verify mock expectations
			mockCtrl.AssertExpectations(t)
		})
	}
}
