package users

import (
	"errors"
	"net/http"

	"omg/api/internal/controller/users"
	"omg/api/internal/model"

	"github.com/gin-gonic/gin"
)

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type registerResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input, err := validateAndMapRegisterRequest(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserAlreadyExists):
			c.JSON(http.StatusBadRequest, gin.H{"error": "user already exists"})
		case errors.Is(err, users.ErrHashedPassword):
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, registerResponse{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status.String(),
	})
}

func validateAndMapRegisterRequest(req registerRequest) (model.CreateUserInput, error) {
	if req.Name == "" {
		return model.CreateUserInput{}, errors.New("user name is required")
	}

	if req.Email == "" {
		return model.CreateUserInput{}, errors.New("user email is required")
	}

	if req.Password == "" {
		return model.CreateUserInput{}, errors.New("user password is required")
	}

	if !emailRegex.MatchString(req.Email) {
		return model.CreateUserInput{}, errors.New("invalid email format")
	}

	return model.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
