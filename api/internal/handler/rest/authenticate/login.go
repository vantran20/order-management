package authenticate

import (
	"errors"
	"net/http"

	"omg/api/internal/authenticate"

	"github.com/gin-gonic/gin"
)

// loginRequest presents the login request fields
type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// loginResponse presents the login response data
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.authService.Login(c, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, authenticate.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Get user ID from the token claims
	claims, err := h.authService.ValidateToken(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate token"})
		return
	}

	refreshToken, err := h.authService.GenerateRefreshToken(claims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
