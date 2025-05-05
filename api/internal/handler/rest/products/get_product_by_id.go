package products

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/products"
	"omg/api/pkg/floatutil"

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

	productID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}

	if productID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	p, err := h.controller.GetByID(c.Request.Context(), productID)
	if err != nil {
		switch {
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, getProductByIDResponse{
		ID:          strconv.FormatInt(p.ID, 10),
		Name:        p.Name,
		Description: p.Description,
		Price:       floatutil.FormatFloat(p.Price),
		Stock:       strconv.FormatInt(p.Stock, 10),
		Status:      p.Status.String(),
	})
}
