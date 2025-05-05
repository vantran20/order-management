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

func TestHandler_UpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockUpdateCtrl struct {
		wantCall bool
		inp      model.UpdateProductInput
		out      model.Product
		err      error
	}

	type arg struct {
		request        updateProductRequest
		mockUpdateCtrl mockUpdateCtrl
		expectedStatus int
		expectedBody   interface{}
	}

	tcs := map[string]arg{
		"successful_update": {
			request: updateProductRequest{
				ID:          "123",
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateProductInput{
					ID:          123,
					Name:        "Test product",
					Description: "Test description",
					Price:       2000,
					Stock:       100,
					Status:      model.ProductStatusActive,
				},
				out: model.Product{
					ID:          123,
					Name:        "Test product",
					Description: "Test description",
					Price:       2000,
					Stock:       100,
					Status:      model.ProductStatusActive,
				},
				err: nil,
			},
			expectedStatus: http.StatusCreated,
			expectedBody: updateProductRequest{
				ID:          "123",
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      model.ProductStatusActive.String(),
			},
		},
		"invalid_request_missing_id": {
			request: updateProductRequest{
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "product id is required"},
		},
		"invalid_request_invalid_id_format": {
			request: updateProductRequest{
				ID:          "abc",
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "strconv.ParseInt: parsing \"abc\": invalid syntax"},
		},
		"invalid_request_missing_name": {
			request: updateProductRequest{
				ID:          "123",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "product name is required"},
		},
		"invalid_request_missing_desc": {
			request: updateProductRequest{
				ID:     "123",
				Name:   "Test product",
				Price:  "2000",
				Stock:  "100",
				Status: "ACTIVE",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "product description is required"},
		},
		"invalid_request_missing_price": {
			request: updateProductRequest{
				ID:          "123",
				Name:        "Test product",
				Description: "Test description",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "product price is required"},
		},
		"invalid_request_invalid_status": {
			request: updateProductRequest{
				ID:          "123",
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "invalid_status",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid product status"},
		},
		"product_not_found": {
			request: updateProductRequest{
				ID:          "123",
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateProductInput{
					ID:          123,
					Name:        "Test product",
					Description: "Test description",
					Price:       2000,
					Stock:       100,
					Status:      model.ProductStatusActive,
				},
				err: products.ErrNotFound,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "product not found"},
		},
		"product_deleted": {
			request: updateProductRequest{
				ID:          "123",
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateProductInput{
					ID:          123,
					Name:        "Test product",
					Description: "Test description",
					Price:       2000,
					Stock:       100,
					Status:      model.ProductStatusActive,
				},
				err: products.ErrProductDeleted,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "product deleted"},
		},
		"internal_server_error": {
			request: updateProductRequest{
				ID:          "123",
				Name:        "Test product",
				Description: "Test description",
				Price:       "2000",
				Stock:       "100",
				Status:      "ACTIVE",
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateProductInput{
					ID:          123,
					Name:        "Test product",
					Description: "Test description",
					Price:       2000,
					Stock:       100,
					Status:      model.ProductStatusActive,
				},
				err: errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "internal server error"},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Setup
			mockCtrl := products.NewMockController(t)
			handler := New(mockCtrl)

			// Create a test router
			router := gin.New()
			router.PUT("/authenticated/products/update", handler.UpdateProduct)

			// Setup mock expectations
			if tc.mockUpdateCtrl.wantCall {
				mockCtrl.On("Update", mock.Anything, tc.mockUpdateCtrl.inp).Return(tc.mockUpdateCtrl.out, tc.mockUpdateCtrl.err)
			}

			// Create test request
			reqBody, _ := json.Marshal(tc.request)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/authenticated/products/update", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			require.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
