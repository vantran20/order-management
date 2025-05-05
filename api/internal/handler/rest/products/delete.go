package products

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/products"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Delete(c *gin.Context) {
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

	if err = h.controller.Delete(c.Request.Context(), productID); err != nil {
		switch {
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, "Delete product successfully")
}
