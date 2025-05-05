package products

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"omg/api/internal/controller/products"
	"omg/api/internal/model"
	"omg/api/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockProductCtrl struct {
		wantCall bool
		input    model.CreateProductInput
		output   model.Product
		err      error
	}
	tests := map[string]struct {
		requestBody     createRequest
		mockProductCtrl mockProductCtrl
		expStatus       int
		expectedBody    interface{}
	}{
		"successful order creation": {
			requestBody: createRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       "2000",
				Stock:       "100",
			},
			mockProductCtrl: mockProductCtrl{
				wantCall: true,
				input: model.CreateProductInput{
					Name:  "Test Product",
					Desc:  "Test Description",
					Price: 2000,
					Stock: 100,
				},
				output: model.Product{
					ID:          1,
					Name:        "Test Product",
					Description: "Test Description",
					Price:       2000,
					Stock:       100,
					Status:      model.ProductStatusActive,
				},
			},
			expStatus: http.StatusCreated,
			expectedBody: createResponse{
				ID:          "1",
				Name:        "Test Product",
				Description: "Test Description",
				Price:       "2000",
				Stock:       "100",
				Status:      model.ProductStatusActive.String(),
			},
		},
		"invalid_name_format": {
			requestBody: createRequest{
				Name:        "",
				Description: "Test Description",
				Price:       "2000",
				Stock:       "100",
			},
			expStatus:    http.StatusBadRequest,
			expectedBody: gin.H{"error": "Key: 'createRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
		"invalid_description_format": {
			requestBody: createRequest{
				Name:        "Test Product",
				Description: "",
				Price:       "2000",
				Stock:       "100",
			},
			expStatus:    http.StatusBadRequest,
			expectedBody: gin.H{"error": "Key: 'createRequest.Description' Error:Field validation for 'Description' failed on the 'required' tag"},
		},
		"invalid_price_format": {
			requestBody: createRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       "",
				Stock:       "100",
			},
			expStatus:    http.StatusBadRequest,
			expectedBody: gin.H{"error": "Key: 'createRequest.Price' Error:Field validation for 'Price' failed on the 'required' tag"},
		},
		"invalid_stock_format": {
			requestBody: createRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       "2000",
				Stock:       "",
			},
			expStatus:    http.StatusBadRequest,
			expectedBody: gin.H{"error": "Key: 'createRequest.Stock' Error:Field validation for 'Stock' failed on the 'required' tag"},
		},
		"product_exists": {
			requestBody: createRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       "2000",
				Stock:       "100",
			},
			mockProductCtrl: mockProductCtrl{
				wantCall: true,
				input: model.CreateProductInput{
					Name:  "Test Product",
					Desc:  "Test Description",
					Price: 2000,
					Stock: 100,
				},
				err: products.ErrProductAlreadyExists,
			},
			expStatus:    http.StatusBadRequest,
			expectedBody: gin.H{"error": "product already exists"},
		},
		"internal server error": {
			requestBody: createRequest{
				Name:        "Test Product",
				Description: "Test Description",
				Price:       "2000",
				Stock:       "100",
			},
			mockProductCtrl: mockProductCtrl{
				wantCall: true,
				input: model.CreateProductInput{
					Name:  "Test Product",
					Desc:  "Test Description",
					Price: 2000,
					Stock: 100,
				},
				err: errors.New("unexpected error"),
			},
			expStatus:    http.StatusInternalServerError,
			expectedBody: gin.H{"error": "internal server error"},
		},
	}

	for desc, tc := range tests {
		t.Run(desc, func(t *testing.T) {
			// Setup
			mockCtrl := products.NewMockController(t)
			handler := New(mockCtrl)

			// Create a test router
			router := gin.New()
			router.POST("/authenticated/products/create", handler.Create)

			// Setup mock expectations
			if tc.mockProductCtrl.wantCall {
				mockCtrl.On("Create", mock.Anything, tc.mockProductCtrl.input).Return(tc.mockProductCtrl.output, tc.mockProductCtrl.err)
			}

			// Create test request
			reqBody, _ := json.Marshal(tc.requestBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/authenticated/products/create", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, tc.expStatus, w.Code)
			require.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
