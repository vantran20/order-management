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

func TestHandler_GetUserByEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type mockGetByEmailCtrl struct {
		wantCall bool
		email    string
		out      model.User
		err      error
	}

	type arg struct {
		request        interface{}
		mockGetByEmail mockGetByEmailCtrl
		expectedStatus int
		expectedBody   interface{}
	}

	tcs := map[string]arg{
		"successful_retrieval": {
			request: getUserByEmailRequest{
				Email: "test@example.com",
			},
			mockGetByEmail: mockGetByEmailCtrl{
				wantCall: true,
				email:    "test@example.com",
				out: model.User{
					ID:     123,
					Name:   "Test User",
					Email:  "test@example.com",
					Status: model.UserStatusActive,
				},
				err: nil,
			},
			expectedStatus: http.StatusCreated,
			expectedBody: getUserByEmailResponse{
				ID:     "123",
				Name:   "Test User",
				Email:  "test@example.com",
				Status: model.UserStatusActive.String(),
			},
		},
		"missing_email": {
			request: getUserByEmailRequest{
				Email: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "email is required"},
		},
		"invalid_email_format": {
			request: getUserByEmailRequest{
				Email: "invalid-email",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid email format"},
		},
		"invalid_email_no_at": {
			request: getUserByEmailRequest{
				Email: "testexample.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid email format"},
		},
		"invalid_email_no_domain": {
			request: getUserByEmailRequest{
				Email: "test@",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid email format"},
		},
		"invalid_email_special_chars": {
			request: getUserByEmailRequest{
				Email: "test!@example.com",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "invalid email format"},
		},
		"user_not_found": {
			request: getUserByEmailRequest{
				Email: "nonexistent@example.com",
			},
			mockGetByEmail: mockGetByEmailCtrl{
				wantCall: true,
				email:    "nonexistent@example.com",
				err:      users.ErrUserNotFound,
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   gin.H{"error": "user not found"},
		},
		"internal_server_error": {
			request: getUserByEmailRequest{
				Email: "test@example.com",
			},
			mockGetByEmail: mockGetByEmailCtrl{
				wantCall: true,
				email:    "test@example.com",
				err:      errors.New("database error"),
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "internal server error"},
		},
		"user_deleted": {
			request: getUserByEmailRequest{
				Email: "deleted@example.com",
			},
			mockGetByEmail: mockGetByEmailCtrl{
				wantCall: true,
				email:    "deleted@example.com",
				out: model.User{
					ID:     456,
					Name:   "Deleted User",
					Email:  "deleted@example.com",
					Status: model.UserStatusDeleted,
				},
				err: nil,
			},
			expectedStatus: http.StatusCreated,
			expectedBody: getUserByEmailResponse{
				ID:     "456",
				Name:   "Deleted User",
				Email:  "deleted@example.com",
				Status: model.UserStatusDeleted.String(),
			},
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Setup
			mockCtrl := users.NewMockController(t)
			handler := New(mockCtrl)

			// Create a test router
			router := gin.New()
			router.GET("/authenticated/users/profile", handler.GetUserByEmail)

			// Setup mock expectations
			if tc.mockGetByEmail.wantCall {
				mockCtrl.On("GetByEmail", mock.Anything, tc.mockGetByEmail.email).Return(tc.mockGetByEmail.out, tc.mockGetByEmail.err)
			}

			// Create test request
			reqBody, _ := json.Marshal(tc.request)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/authenticated/users/profile", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.JSONEq(t, testutil.ToJSONString(tc.expectedBody), w.Body.String())
			mockCtrl.AssertExpectations(t)
		})
	}
}

func TestValidateRequest(t *testing.T) {
	tcs := map[string]struct {
		input    getUserByEmailRequest
		expected string
		err      error
	}{
		"valid_email": {
			input: getUserByEmailRequest{
				Email: "test@example.com",
			},
			expected: "test@example.com",
			err:      nil,
		},
		"valid_email_with_subdomain": {
			input: getUserByEmailRequest{
				Email: "test@sub.example.com",
			},
			expected: "test@sub.example.com",
			err:      nil,
		},
		"valid_email_with_dots": {
			input: getUserByEmailRequest{
				Email: "test.user@example.com",
			},
			expected: "test.user@example.com",
			err:      nil,
		},
		"valid_email_with_plus": {
			input: getUserByEmailRequest{
				Email: "test+user@example.com",
			},
			expected: "test+user@example.com",
			err:      nil,
		},
		"valid_email_with_underscore": {
			input: getUserByEmailRequest{
				Email: "test_user@example.com",
			},
			expected: "test_user@example.com",
			err:      nil,
		},
		"valid_email_with_numbers": {
			input: getUserByEmailRequest{
				Email: "test123@example.com",
			},
			expected: "test123@example.com",
			err:      nil,
		},
		"empty_email": {
			input: getUserByEmailRequest{
				Email: "",
			},
			expected: "",
			err:      errors.New("email is required"),
		},
		"missing_at": {
			input: getUserByEmailRequest{
				Email: "testexample.com",
			},
			expected: "",
			err:      errors.New("invalid email format"),
		},
		"missing_domain": {
			input: getUserByEmailRequest{
				Email: "test@",
			},
			expected: "",
			err:      errors.New("invalid email format"),
		},
		"invalid_special_chars": {
			input: getUserByEmailRequest{
				Email: "test!@example.com",
			},
			expected: "",
			err:      errors.New("invalid email format"),
		},
		"invalid_domain_chars": {
			input: getUserByEmailRequest{
				Email: "test@example!.com",
			},
			expected: "",
			err:      errors.New("invalid email format"),
		},
		"invalid_tld": {
			input: getUserByEmailRequest{
				Email: "test@example.c",
			},
			expected: "",
			err:      errors.New("invalid email format"),
		},
		"multiple_at": {
			input: getUserByEmailRequest{
				Email: "test@example@com",
			},
			expected: "",
			err:      errors.New("invalid email format"),
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			email, err := validateRequest(tc.input)
			if tc.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, email)
		})
	}
}
