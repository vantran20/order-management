package products

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"omg/api/internal/controller/products"
	"omg/api/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockDeleteProductCtrl struct {
		wantCall bool
		id       int64
		err      error
	}
	tests := map[string]struct {
		givenID               string
		mockDeleteProductCtrl mockDeleteProductCtrl
		expStatus             int
		expectedBody          interface{}
	}{
		"successful order creation": {
			givenID: "1",
			mockDeleteProductCtrl: mockDeleteProductCtrl{
				wantCall: true,
				id:       1,
			},
			expStatus:    http.StatusOK,
			expectedBody: "Delete product successfully",
		},
		"invalid_product_id_format": {
			givenID:      "abc",
			expStatus:    http.StatusBadRequest,
			expectedBody: gin.H{"error": "invalid product ID format"},
		},
		"invalid_product_id": {
			givenID:      "0",
			expStatus:    http.StatusBadRequest,
			expectedBody: gin.H{"error": "invalid product ID"},
		},
		"product_not_found": {
			givenID: "1",
			mockDeleteProductCtrl: mockDeleteProductCtrl{
				wantCall: true,
				id:       1,
				err:      products.ErrNotFound,
			},
			expStatus:    http.StatusNotFound,
			expectedBody: gin.H{"error": "product not found"},
		},
		"internal server error": {
			givenID: "1",
			mockDeleteProductCtrl: mockDeleteProductCtrl{
				wantCall: true,
				id:       1,
				err:      errors.New("unexpected error"),
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
			router.POST("/authenticated/products/delete/:id", handler.Delete)

			// Setup mock expectations
			if tc.mockDeleteProductCtrl.wantCall {
				mockCtrl.On("Delete", mock.Anything, tc.mockDeleteProductCtrl.id).Return(tc.mockDeleteProductCtrl.err)
			}

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/authenticated/products/delete/"+tc.givenID, nil)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, tc.expStatus, w.Code)
			require.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
