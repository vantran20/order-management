package users

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"omg/api/internal/controller/users"

	"github.com/gin-gonic/gin"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type getUserByEmailRequest struct {
	Email string `json:"email"`
}

type getUserByEmailResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// GetUserByEmail handles user registration
func (h *Handler) GetUserByEmail(c *gin.Context) {
	var req getUserByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := validateRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.controller.GetByEmail(c.Request.Context(), email)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			c.JSON(http.StatusConflict, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, getUserByEmailResponse{
		ID:     strconv.FormatInt(user.ID, 10),
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status.String(),
	})
}

func validateRequest(inp getUserByEmailRequest) (string, error) {
	if inp.Email == "" {
		return "", errors.New("email is required")
	}

	if !emailRegex.MatchString(inp.Email) {
		return "", errors.New("invalid email format")
	}

	return inp.Email, nil
}
