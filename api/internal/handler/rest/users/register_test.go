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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
		})
	}
}
