package products

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type getProductsResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
	Status      string `json:"status"`
}

func (h *Handler) List(c *gin.Context) {
	list, err := h.controller.List(c.Request.Context())
	if err != nil {
		switch {
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	var response []getProductsResponse
	for _, p := range list {
		response = append(response, getProductsResponse{
			ID:          strconv.FormatInt(p.ID, 10),
			Name:        p.Name,
			Description: p.Description,
			Price:       strconv.FormatFloat(p.Price, 'f', -1, 64),
			Stock:       strconv.FormatInt(p.Stock, 10),
			Status:      p.Status.String(),
		})
	}

	c.JSON(http.StatusOK, response)
}
