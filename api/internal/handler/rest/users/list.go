package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getUserResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

func (h *Handler) List(c *gin.Context) {
	list, err := h.controller.GetUsers(c.Request.Context())
	if err != nil {
		switch {
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	var response []getUserResponse
	for _, user := range list {
		response = append(response, getUserResponse{
			ID:     user.ID,
			Name:   user.Name,
			Email:  user.Email,
			Status: user.Status.String(),
		})
	}

	c.JSON(http.StatusOK, response)
}
