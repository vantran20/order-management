package authenticate

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestMiddlewareAuth_Handler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type arg struct {
		name           string
		authHeader     string
		mockClaims     *Claims
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
		expectedUserID int64
		expectedEmail  string
	}

	tcs := map[string]arg{
		"missing authorization header": {
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header is required",
			},
		},
		"invalid header format": {
			name:           "invalid header format",
			authHeader:     "InvalidFormat token123",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid authorization header format",
			},
		},
		"expired token": {
			name: "expired token",
			mockClaims: &Claims{
				UserID: 123,
				Email:  "test@example.com",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)), // Expired 24 hours ago
					IssuedAt:  jwt.NewNumericDate(time.Now().Add(-48 * time.Hour)), // Issued 48 hours ago
				},
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Token expired",
			},
		},
		"invalid token": {
			name:           "invalid token",
			authHeader:     "Bearer invalid.token",
			mockError:      errors.New("invalid token"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid token",
			},
		},
		"successful authentication": {
			name:       "successful authentication",
			authHeader: "Bearer valid.token",
			mockClaims: &Claims{
				UserID: 123,
				Email:  "test@example.com",
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
			expectedStatus: http.StatusOK,
			expectedUserID: 123,
			expectedEmail:  "test@example.com",
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Setup
			testSecret := []byte("test-secret-key")
			authService := AuthService{
				secret: testSecret,
			}

			// Create a valid token for successful case
			if tc.mockClaims != nil {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, tc.mockClaims)
				signedToken, err := token.SignedString(testSecret)
				require.NoError(t, err)
				tc.authHeader = "Bearer " + signedToken
			}

			middleware := &MiddlewareAuth{
				authService: authService,
			}

			// Create test request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/", nil)
			if tc.authHeader != "" {
				c.Request.Header.Set("Authorization", tc.authHeader)
			}

			// Create a channel to capture context values
			userIDChan := make(chan int64, 1)
			emailChan := make(chan string, 1)

			// Add test handlers
			handlers := []gin.HandlerFunc{
				middleware.Handler(),
				func(c *gin.Context) {
					if val, exists := c.Get("user_id"); exists {
						userIDChan <- val.(int64)
					}
					if val, exists := c.Get("email"); exists {
						emailChan <- val.(string)
					}
					c.Status(http.StatusOK)
				},
			}

			// Execute handlers
			for _, handler := range handlers {
				handler(c)
				if c.IsAborted() {
					break
				}
			}

			// Assert
			require.Equal(t, tc.expectedStatus, w.Code)
			if tc.expectedBody != nil {
				require.JSONEq(t, toJSON(tc.expectedBody), w.Body.String())
			}

			if tc.expectedStatus == http.StatusOK {
				select {
				case userID := <-userIDChan:
					require.Equal(t, tc.expectedUserID, userID)
				default:
					t.Error("Expected user_id not found in context")
				}
				select {
				case email := <-emailChan:
					require.Equal(t, tc.expectedEmail, email)
				default:
					t.Error("Expected email not found in context")
				}
			}
		})
	}
}

// Helper function to convert map to JSON string
func toJSON(m map[string]interface{}) string {
	jsonStr := "{"
	for k, v := range m {
		jsonStr += `"` + k + `":"` + v.(string) + `",`
	}
	jsonStr = jsonStr[:len(jsonStr)-1] + "}"
	return jsonStr
}
