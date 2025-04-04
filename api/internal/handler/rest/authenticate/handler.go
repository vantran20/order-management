package authenticate

import (
	"omg/api/internal/authenticate"
)

type Handler struct {
	authService authenticate.AuthService
}

func New(authService authenticate.AuthService) Handler {
	return Handler{
		authService: authService,
	}
}
