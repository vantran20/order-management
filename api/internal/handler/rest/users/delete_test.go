package users

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"omg/api/internal/controller/users"
	"omg/api/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockDeleteCtrl struct {
		wantCall bool
		id       int64
		err      error
	}

	type arg struct {
		userID         string
		mockDeleteCtrl mockDeleteCtrl
		expectedStatus int
		expectedBody   interface{}
	}

	tcs := map[string]arg{
		"successful_deletion": {
			userID: "123",
			mockDeleteCtrl: mockDeleteCtrl{
				wantCall: true,
				id:       123,
				err:      nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "Delete user successfully",
		},
		"invalid_user_id_format": {
			userID:         "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid user ID format"},
		},
		"user_not_found": {
			userID: "123",
			mockDeleteCtrl: mockDeleteCtrl{
				wantCall: true,
				id:       123,
				err:      users.ErrUserNotFound,
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "user not found"},
		},
		"internal_server_error": {
			userID: "123",
			mockDeleteCtrl: mockDeleteCtrl{
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
			mockCtrl := users.NewMockController(t)
			handler := New(mockCtrl)

			// Create a test router
			router := gin.New()
			router.POST("/users/delete/:id", handler.Delete)

			// Setup mock expectations
			if tc.mockDeleteCtrl.wantCall {
				mockCtrl.On("Delete", mock.Anything, tc.mockDeleteCtrl.id).Return(tc.mockDeleteCtrl.err)
			}

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/users/delete/"+tc.userID, nil)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
