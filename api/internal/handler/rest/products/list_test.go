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

func TestHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockGetListCtrl struct {
		wantCall bool
		out      []model.Product
		err      error
	}

	type arg struct {
		mockGetListCtrl mockGetListCtrl
		expectedStatus  int
		expectedBody    interface{}
	}

	tcs := map[string]arg{
		"successful_retrieval": {
			mockGetListCtrl: mockGetListCtrl{
				wantCall: true,
				out: []model.Product{
					{
						ID:          123,
						Name:        "test product",
						Description: "test description",
						Price:       2000,
						Stock:       100,
						Status:      model.ProductStatusActive,
					},
				},
				err: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody: []getProductByIDResponse{
				{
					ID:          "123",
					Name:        "test product",
					Description: "test description",
					Price:       "2000",
					Stock:       "100",
					Status:      model.ProductStatusActive.String(),
				},
			},
		},
		"empty_retrieval": {
			mockGetListCtrl: mockGetListCtrl{
				wantCall: true,
				out:      nil,
				err:      nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		"internal_server_error": {
			mockGetListCtrl: mockGetListCtrl{
				wantCall: true,
				out: []model.Product{
					{
						ID:          123,
						Name:        "test product",
						Description: "test description",
						Price:       2000,
						Stock:       100,
						Status:      model.ProductStatusActive,
					},
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
			router.GET("/authenticated/products/list", handler.List)

			// Setup mock expectations
			if tc.mockGetListCtrl.wantCall {
				mockCtrl.On("List", mock.Anything).Return(tc.mockGetListCtrl.out, tc.mockGetListCtrl.err)
			}

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/authenticated/products/list", nil)
			router.ServeHTTP(w, req)

			// Assertions
			require.Equal(t, tc.expectedStatus, w.Code)
			require.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
