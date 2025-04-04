package products

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/users"

	"github.com/gin-gonic/gin"
)

type getProductByIDResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
	Status      string `json:"status"`
}

func (h *Handler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	p, err := h.controller.GetByID(c.Request.Context(), productID)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, getProductByIDResponse{
		ID:          strconv.FormatInt(p.ID, 10),
		Name:        p.Name,
		Description: p.Description,
		Price:       strconv.FormatFloat(p.Price, 'f', -1, 64),
		Stock:       strconv.FormatInt(p.Stock, 10),
		Status:      p.Status.String(),
	})
}
