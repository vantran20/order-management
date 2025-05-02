package authenticate

import (
	"errors"
	"math"
	"net/http/httptest"
	"testing"
	"time"

	"omg/api/internal/model"
	"omg/api/internal/repository"
	"omg/api/internal/repository/user"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	// Setup test claims secret for jwt token generation
	testSecret := []byte("test-secret-key")

	type arg struct {
		givenEmail     string
		givenPassword  string
		mockUser       model.User
		mockUserErr    error
		mockBcryptErr  error
		mockSigningErr error
		expToken       string
		expErr         error
	}

	tcs := map[string]arg{
		"success": {
			givenEmail:    "test@example.com",
			givenPassword: "password123",
			mockUser: model.User{
				ID:       14753001,
				Email:    "test@example.com",
				Password: "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", // mock bcrypt hash
			},
			expToken: "valid.jwt.token", // This will be replaced with actual token in the test
		},
		"user_not_found": {
			givenEmail:    "notfound@example.com",
			givenPassword: "password123",
			mockUserErr:   user.ErrNotFound,
			expErr:        ErrInvalidCredentials,
		},
		"database_error": {
			givenEmail:    "test@example.com",
			givenPassword: "password123",
			mockUserErr:   errors.New("database error"),
			expErr:        errors.New("database error"),
		},
		"invalid_password": {
			givenEmail:    "test@example.com",
			givenPassword: "wrongpassword",
			mockUser: model.User{
				ID:       14753001,
				Email:    "test@example.com",
				Password: "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", // mock bcrypt hash
			},
			mockBcryptErr: bcrypt.ErrMismatchedHashAndPassword,
			expErr:        ErrInvalidCredentials,
		},
		"jwt_signing_error": {
			givenEmail:    "test@example.com",
			givenPassword: "password123",
			mockUser: model.User{
				ID:       14753001,
				Email:    "test@example.com",
				Password: "$2a$10$XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", // mock bcrypt hash
			},
			mockSigningErr: errors.New("signing error"),
			expErr:         errors.New("signing error"),
		},
	}

	for s, tc := range tcs {
		t.Run(s, func(t *testing.T) {
			// Given:
			// Create a mock gin context
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Setup mocks
			mockUserRepo := &user.MockRepository{}
			mockRepo := &repository.MockRegistry{}

			// Setup mock user repository
			if tc.givenEmail != "" {
				mockUserRepo.On("GetByEmail", mock.Anything, tc.givenEmail).Return(tc.mockUser, tc.mockUserErr)
			}
			mockRepo.On("User").Return(mockUserRepo)

			// Mock bcrypt.CompareHashAndPassword
			mockCompareHashAndPassword := func(hashedPassword, password []byte) error {
				if tc.mockBcryptErr != nil {
					return tc.mockBcryptErr
				}
				return nil
			}

			// Create service with mocked dependencies
			authService := &AuthService{
				repo:                   mockRepo,
				secret:                 testSecret,
				compareHashAndPassword: mockCompareHashAndPassword,
				signToken: func(token *jwt.Token, secret []byte) (string, error) {
					return token.SignedString(secret)
				},
			}

			// Mock jwt.Token.SignedString for error case
			if tc.mockSigningErr != nil {
				mockSignToken := func(token *jwt.Token, secret []byte) (string, error) {
					return "", tc.mockSigningErr
				}
				authService.signToken = mockSignToken
			}

			// When:
			token, err := authService.Login(c, tc.givenEmail, tc.givenPassword)

			// Then:
			mockUserRepo.AssertExpectations(t)
			mockRepo.AssertExpectations(t)

			if tc.expErr != nil {
				require.Error(t, err)
				if errors.Is(err, tc.expErr) {
					require.True(t, errors.Is(err, tc.expErr))
				} else {
					require.Equal(t, tc.expErr.Error(), err.Error())
				}
				require.Empty(t, token)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, token)

				// Verify JWT token
				parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
					return testSecret, nil
				})
				require.NoError(t, err)
				require.True(t, parsedToken.Valid)

				// Verify claims
				claims, ok := parsedToken.Claims.(*Claims)
				require.True(t, ok)
				require.Equal(t, tc.mockUser.ID, claims.UserID)
				require.Equal(t, tc.mockUser.Email, claims.Email)
				require.NotZero(t, claims.ExpiresAt)
				require.NotZero(t, claims.IssuedAt)

				// Check that expiry is around 24 hours in the future (with 5 minute margin for test execution)
				expectedExpiry := time.Now().Add(24 * time.Hour)
				tokenExpiry := claims.ExpiresAt.Time
				timeDiff := expectedExpiry.Sub(tokenExpiry)
				require.LessOrEqual(t, math.Abs(timeDiff.Minutes()), float64(5))
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	testSecret := []byte("test-secret-key")

	type arg struct {
		token         string
		expectedUser  int64
		expectedEmail string
		expectedError error
	}

	validClaims := Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	validToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims).SignedString(testSecret)
	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-48 * time.Hour)),
		},
	}).SignedString(testSecret)

	tcs := map[string]arg{
		"valid token": {
			token:         validToken,
			expectedUser:  1,
			expectedEmail: "test@example.com",
		},
		"expired token": {
			token:         expiredToken,
			expectedError: ErrTokenExpired,
		},
	}

	for desc, tc := range tcs {
		t.Run(desc, func(t *testing.T) {
			// Setup
			authService := &AuthService{
				secret: testSecret,
			}

			// Execute
			claims, err := authService.ValidateToken(tc.token)

			// Assert
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError.Error(), err.Error())
				require.Nil(t, claims)
			} else {
				require.NoError(t, err)
				require.NotNil(t, claims)
				require.Equal(t, tc.expectedUser, claims.UserID)
				require.Equal(t, tc.expectedEmail, claims.Email)
			}
		})
	}
}

func TestAuthService_GenerateRefreshToken(t *testing.T) {
	testSecret := []byte("test-secret-key")
	userID := int64(1)

	// Setup
	authService := &AuthService{
		secret: testSecret,
	}

	// Execute
	token, err := authService.GenerateRefreshToken(userID)

	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify token
	claims := jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return testSecret, nil
	})
	require.NoError(t, err)
	require.True(t, parsedToken.Valid)
	require.Equal(t, "1", claims.ID)
	require.NotZero(t, claims.ExpiresAt)
	require.NotZero(t, claims.IssuedAt)
}
