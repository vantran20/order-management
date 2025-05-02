package authenticate

import (
	"omg/api/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth interface {
	Login(ctx *gin.Context, email, password string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
	GenerateRefreshToken(userID int64) (string, error)
}

func NewAuthService(repo repository.Registry, secret string) AuthService {
	return AuthService{
		repo:                   repo,
		secret:                 []byte(secret),
		compareHashAndPassword: bcrypt.CompareHashAndPassword,
		signToken: func(token *jwt.Token, secret []byte) (string, error) {
			return token.SignedString(secret)
		},
	}
}

type AuthService struct {
	repo                   repository.Registry
	secret                 []byte
	compareHashAndPassword func(hashedPassword, password []byte) error
	signToken              func(token *jwt.Token, secret []byte) (string, error)
}

type Middleware interface {
	Handler() gin.HandlerFunc
}

func NewAuthMiddleware(authService AuthService) *MiddlewareAuth {
	return &MiddlewareAuth{
		authService: authService,
	}
}

type MiddlewareAuth struct {
	authService AuthService
}
