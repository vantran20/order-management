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
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Status   string `json:"status" binding:"required"`
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

	input, err := validateAndMapRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.controller.Update(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
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

func validateAndMapRequest(req updateUserRequest) (model.UpdateUserInput, error) {
	if req.ID == "" {
		return model.UpdateUserInput{}, errors.New("user ID is required")
	}

	if req.Name == "" {
		return model.UpdateUserInput{}, errors.New("user name is required")
	}

	if req.Email == "" {
		return model.UpdateUserInput{}, errors.New("user email is required")
	}

	if !emailRegex.MatchString(req.Email) {
		return model.UpdateUserInput{}, errors.New("invalid email format")
	}

	if req.Status == "" {
		return model.UpdateUserInput{}, errors.New("user status is required")
	}

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return model.UpdateUserInput{}, errors.New("invalid user ID format")
	}

	status := model.UserStatus(req.Status)
	if !status.IsValid() {
		return model.UpdateUserInput{}, errors.New("invalid user status")
	}

	return model.UpdateUserInput{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Status:   model.UserStatus(req.Status),
	}, nil
}
