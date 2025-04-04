package users

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/users"

	"github.com/gin-gonic/gin"
)

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

	user, err := h.controller.GetByEmail(c.Request.Context(), req.Email)
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
