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

func TestHandler_Register(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockCreateCtrl struct {
		wantCall bool
		inp      model.CreateUserInput
		out      model.User
		err      error
	}

	type arg struct {
		request        registerRequest
		mockCreateCtrl mockCreateCtrl
		expectedStatus int
		expectedBody   interface{}
	}

	tcs := map[string]arg{
		"successful_registration": {
			request: registerRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockCreateCtrl: mockCreateCtrl{
				wantCall: true,
				inp: model.CreateUserInput{
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password123",
				},
				out: model.User{
					ID:     123,
					Name:   "Test User",
					Email:  "test@example.com",
					Status: model.UserStatusActive,
				},
				err: nil,
			},
			expectedStatus: http.StatusCreated,
			expectedBody: registerResponse{
				ID:     123,
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive.String(),
			},
		},
		"invalid_request_missing_name": {
			request: registerRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'registerRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag"},
		},
		"invalid_request_missing_email": {
			request: registerRequest{
				Name:     "Test User",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'registerRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag"},
		},
		"invalid_request_missing_password": {
			request: registerRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'registerRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag"},
		},
		"invalid_request_invalid_email": {
			request: registerRequest{
				Name:     "Test User",
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "Key: 'registerRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"},
		},
		"user_already_exists": {
			request: registerRequest{
				Name:     "Test User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockCreateCtrl: mockCreateCtrl{
				wantCall: true,
				inp: model.CreateUserInput{
					Name:     "Test User",
					Email:    "existing@example.com",
					Password: "password123",
				},
				err: users.ErrUserAlreadyExists,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "user already exists"},
		},
		"password_hashing_error": {
			request: registerRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockCreateCtrl: mockCreateCtrl{
				wantCall: true,
				inp: model.CreateUserInput{
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password123",
				},
				err: users.ErrHashedPassword,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "failed to hash password"},
		},
		"internal_server_error": {
			request: registerRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockCreateCtrl: mockCreateCtrl{
				wantCall: true,
				inp: model.CreateUserInput{
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "password123",
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
			router.POST("/users", handler.Register)

			// Setup mock expectations
			if tc.mockCreateCtrl.wantCall {
				mockCtrl.On("Create", mock.Anything, tc.mockCreateCtrl.inp).Return(tc.mockCreateCtrl.out, tc.mockCreateCtrl.err)
			}

			// Create test request
			reqBody, _ := json.Marshal(tc.request)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assertions
			require.Equal(t, tc.expectedStatus, w.Code)
			require.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}

func Test_validateAndMapRegisterRequest(t *testing.T) {
	tcs := map[string]struct {
		input    registerRequest
		expected model.CreateUserInput
		err      error
	}{
		"valid_request": {
			input: registerRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			expected: model.CreateUserInput{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			err: nil,
		},
		"empty_user_name": {
			input: registerRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			expected: model.CreateUserInput{},
			err:      errors.New("user name is required"),
		},
		"empty_user_email": {
			input: registerRequest{
				Name:     "Test User",
				Email:    "",
				Password: "password123",
			},
			expected: model.CreateUserInput{},
			err:      errors.New("user email is required"),
		},
		"empty_user_pwd": {
			input: registerRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "",
			},
			expected: model.CreateUserInput{},
			err:      errors.New("user password is required"),
		},
		"invalid_user_email": {
			input: registerRequest{
				Name:     "Test User",
				Email:    "invalidtestexample.com",
				Password: "password123",
			},
			expected: model.CreateUserInput{},
			err:      errors.New("invalid email format"),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			output, err := validateAndMapRegisterRequest(tc.input)
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
