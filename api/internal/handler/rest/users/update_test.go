package users

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"omg/api/internal/controller/users"
	"omg/api/internal/model"
	"omg/api/pkg/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_UpdateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockUpdateCtrl struct {
		wantCall bool
		inp      model.UpdateUserInput
		out      model.User
		err      error
	}

	type arg struct {
		request        updateUserRequest
		mockUpdateCtrl mockUpdateCtrl
		expectedStatus int
		expectedBody   interface{}
	}

	tcs := map[string]arg{
		"successful_update": {
			request: updateUserRequest{
				ID:       "123",
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   "ACTIVE",
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateUserInput{
					ID:       123,
					Name:     "Updated User",
					Email:    "updated@example.com",
					Password: "newpassword123",
					Status:   model.UserStatusActive,
				},
				out: model.User{
					ID:     123,
					Name:   "Updated User",
					Email:  "updated@example.com",
					Status: model.UserStatusActive,
				},
				err: nil,
			},
			expectedStatus: http.StatusCreated,
			expectedBody: updateUserResponse{
				ID:     "123",
				Name:   "Updated User",
				Email:  "updated@example.com",
				Status: model.UserStatusActive.String(),
			},
		},
		"invalid_request_missing_id": {
			request: updateUserRequest{
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'updateUserRequest.ID' Error:Field validation for 'ID' failed on the 'required' tag"},
		},
		"invalid_request_invalid_id_format": {
			request: updateUserRequest{
				ID:       "abc",
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid user ID format"},
		},
		"invalid_request_missing_name": {
			request: updateUserRequest{
				ID:       "123",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'updateUserRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
		"invalid_request_missing_email": {
			request: updateUserRequest{
				ID:       "123",
				Name:     "Updated User",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'updateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"},
		},
		"invalid_request_invalid_email": {
			request: updateUserRequest{
				ID:       "123",
				Name:     "Updated User",
				Email:    "invalid-email",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'updateUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"},
		},
		"invalid_request_invalid_status": {
			request: updateUserRequest{
				ID:       "123",
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   "invalid_status",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid user status"},
		},
		"user_not_found": {
			request: updateUserRequest{
				ID:       "123",
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateUserInput{
					ID:       123,
					Name:     "Updated User",
					Email:    "updated@example.com",
					Password: "newpassword123",
					Status:   model.UserStatusActive,
				},
				err: users.ErrUserNotFound,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "user not found"},
		},
		"password_hashing_error": {
			request: updateUserRequest{
				ID:       "123",
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateUserInput{
					ID:       123,
					Name:     "Updated User",
					Email:    "updated@example.com",
					Password: "newpassword123",
					Status:   model.UserStatusActive,
				},
				err: users.ErrHashedPassword,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "failed to hash password"},
		},
		"internal_server_error": {
			request: updateUserRequest{
				ID:       "123",
				Name:     "Updated User",
				Email:    "updated@example.com",
				Password: "newpassword123",
				Status:   model.UserStatusActive.String(),
			},
			mockUpdateCtrl: mockUpdateCtrl{
				wantCall: true,
				inp: model.UpdateUserInput{
					ID:       123,
					Name:     "Updated User",
					Email:    "updated@example.com",
					Password: "newpassword123",
					Status:   model.UserStatusActive,
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
			mockCtrl := users.NewMockController(t)
			handler := New(mockCtrl)

			// Create a test router
			router := gin.New()
			router.PUT("/authenticated/users/update", handler.UpdateUser)

			// Setup mock expectations
			if tc.mockUpdateCtrl.wantCall {
				mockCtrl.On("Update", mock.Anything, tc.mockUpdateCtrl.inp).Return(tc.mockUpdateCtrl.out, tc.mockUpdateCtrl.err)
			}

			// Create test request
			reqBody, _ := json.Marshal(tc.request)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/authenticated/users/update", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			require.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}

func TestValidateAndMapRequest(t *testing.T) {
	tcs := map[string]struct {
		input    updateUserRequest
		expected model.UpdateUserInput
		err      error
	}{
		"valid_request": {
			input: updateUserRequest{
				ID:       "123",
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive.String(),
			},
			expected: model.UpdateUserInput{
				ID:       123,
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive,
			},
			err: nil,
		},
		"invalid_id_format": {
			input: updateUserRequest{
				ID:       "abc",
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive.String(),
			},
			expected: model.UpdateUserInput{},
			err:      errors.New("invalid user ID format"),
		},
		"invalid_user_name": {
			input: updateUserRequest{
				ID:       "abc",
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
				Status:   model.UserStatusActive.String(),
			},
			expected: model.UpdateUserInput{},
			err:      errors.New("user name is required"),
		},
		"invalid_user_email_format": {
			input: updateUserRequest{
				ID:       "abc",
				Name:     "Test User",
				Email:    "invalidtestexample.com",
				Password: "password123",
				Status:   model.UserStatusActive.String(),
			},
			expected: model.UpdateUserInput{},
			err:      errors.New("invalid email format"),
		},
		"invalid_status": {
			input: updateUserRequest{
				ID:       "123",
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Status:   "invalid_status",
			},
			expected: model.UpdateUserInput{},
			err:      errors.New("invalid user status"),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			output, err := validateAndMapRequest(tc.input)
			if tc.err != nil {
				require.Error(t, err)
				require.Equal(t, tc.err.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.expected, output)
		})
	}
}
