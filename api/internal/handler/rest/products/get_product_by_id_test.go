package products

import (
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

func TestHandler_GetProductByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockGetByIDCtrl struct {
		wantCall bool
		id       int64
		out      model.Product
		err      error
	}

	type arg struct {
		givenID         string
		mockGetByIDCtrl mockGetByIDCtrl
		expectedStatus  int
		expectedBody    interface{}
	}

	tcs := map[string]arg{
		"successful_retrieval": {
			givenID: "123",
			mockGetByIDCtrl: mockGetByIDCtrl{
				wantCall: true,
				id:       123,
				out: model.Product{
					ID:          123,
					Name:        "test product",
					Description: "test description",
					Price:       2000,
					Stock:       100,
					Status:      model.ProductStatusActive,
				},
				err: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody: getProductByIDResponse{
				ID:          "123",
				Name:        "test product",
				Description: "test description",
				Price:       "2000",
				Stock:       "100",
				Status:      model.ProductStatusActive.String(),
			},
		},
		"invalid_product_id_format": {
			givenID:        "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid product ID format"},
		},
		"invalid_product_id": {
			givenID:        "0",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid product ID"},
		},
		"product_not_found": {
			givenID: "123",
			mockGetByIDCtrl: mockGetByIDCtrl{
				wantCall: true,
				id:       123,
				err:      products.ErrNotFound,
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "product not found"},
		},
		"internal_server_error": {
			givenID: "123",
			mockGetByIDCtrl: mockGetByIDCtrl{
				wantCall: true,
				id:       123,
				err:      errors.New("database error"),
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
			router.GET("/authenticated/products/:id", handler.GetProductByID)

			// Setup mock expectations
			if tc.mockGetByIDCtrl.wantCall {
				mockCtrl.On("GetByID", mock.Anything, tc.mockGetByIDCtrl.id).Return(tc.mockGetByIDCtrl.out, tc.mockGetByIDCtrl.err)
			}

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/authenticated/products/"+tc.givenID, nil)
			router.ServeHTTP(w, req)

			// Assertions
			require.Equal(t, tc.expectedStatus, w.Code)
			require.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
