package users

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"omg/api/internal/controller/users"
	"omg/api/internal/model"
	"omg/api/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetUserByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockGetByIDCtrl struct {
		wantCall bool
		id       int64
		out      model.User
		err      error
	}

	type arg struct {
		userID          string
		mockGetByIDCtrl mockGetByIDCtrl
		expectedStatus  int
		expectedBody    interface{}
	}

	tcs := map[string]arg{
		"successful_retrieval": {
			userID: "123",
			mockGetByIDCtrl: mockGetByIDCtrl{
				wantCall: true,
				id:       123,
				out: model.User{
					ID:     123,
					Name:   "Test User",
					Email:  "test@example.com",
					Status: model.UserStatusActive,
				},
				err: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody: getUserByIDResponse{
				ID:     "123",
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive.String(),
			},
		},
		"invalid_user_id_format": {
			userID:         "abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid user ID format"},
		},
		"user_not_found": {
			userID: "123",
			mockGetByIDCtrl: mockGetByIDCtrl{
				wantCall: true,
				id:       123,
				err:      users.ErrUserNotFound,
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "user not found"},
		},
		"internal_server_error": {
			userID: "123",
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
			mockCtrl := users.NewMockController(t)
			handler := New(mockCtrl)

			// Create a test router
			router := gin.New()
			router.GET("/authenticated/users/:id", handler.GetUserByID)

			// Setup mock expectations
			if tc.mockGetByIDCtrl.wantCall {
				mockCtrl.On("GetByID", mock.Anything, tc.mockGetByIDCtrl.id).Return(tc.mockGetByIDCtrl.out, tc.mockGetByIDCtrl.err)
			}

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/authenticated/users/"+tc.userID, nil)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
