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

func TestHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockGetUsersCtrl struct {
		wantCall bool
		out      []model.User
		err      error
	}

	type arg struct {
		mockGetUsersCtrl mockGetUsersCtrl
		expectedStatus   int
		expectedBody     interface{}
	}

	tcs := map[string]arg{
		"successful_retrieval": {
			mockGetUsersCtrl: mockGetUsersCtrl{
				wantCall: true,
				out: []model.User{
					{
						ID:     1,
						Name:   "User 1",
						Email:  "user1@example.com",
						Status: model.UserStatusActive,
					},
					{
						ID:     2,
						Name:   "User 2",
						Email:  "user2@example.com",
						Status: model.UserStatusDeleted,
					},
				},
				err: nil,
			},
			expectedStatus: http.StatusOK,
			expectedBody: []getUserResponse{
				{
					ID:     1,
					Name:   "User 1",
					Email:  "user1@example.com",
					Status: model.UserStatusActive.String(),
				},
				{
					ID:     2,
					Name:   "User 2",
					Email:  "user2@example.com",
					Status: model.UserStatusDeleted.String(),
				},
			},
		},
		"empty_list": {
			mockGetUsersCtrl: mockGetUsersCtrl{
				wantCall: true,
				out:      []model.User{},
				err:      nil,
			},
			expectedStatus: http.StatusOK,
		},
		"internal_server_error": {
			mockGetUsersCtrl: mockGetUsersCtrl{
				wantCall: true,
				out:      nil,
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
			router.GET("/authenticated/users", handler.List)

			// Setup mock expectations
			if tc.mockGetUsersCtrl.wantCall {
				mockCtrl.On("GetUsers", mock.Anything).Return(tc.mockGetUsersCtrl.out, tc.mockGetUsersCtrl.err)
			}

			// Create test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/authenticated/users", nil)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
