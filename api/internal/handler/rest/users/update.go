package users

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/users"
	"omg/api/internal/model"

	"github.com/gin-gonic/gin"
)

type updateUserRequest struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Status   string `json:"status"`
}

type updateUserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// UpdateUser handles user update
func (h *Handler) UpdateUser(c *gin.Context) {
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	input := model.UpdateUserInput{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Status:   model.UserStatus(req.Status),
	}

	user, err := h.controller.Update(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			c.JSON(http.StatusConflict, gin.H{"error": "user not found"})
		case errors.Is(err, users.ErrHashedPassword):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, updateUserResponse{
		ID:     strconv.FormatInt(user.ID, 10),
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status.String(),
	})
}
