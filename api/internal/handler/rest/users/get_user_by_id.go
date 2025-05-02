package users

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/users"

	"github.com/gin-gonic/gin"
)

type getUserByIDResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (h *Handler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	user, err := h.controller.GetByID(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, getUserByIDResponse{
		ID:     strconv.FormatInt(user.ID, 10),
		Name:   user.Name,
		Email:  user.Email,
		Status: user.Status.String(),
	})
}
