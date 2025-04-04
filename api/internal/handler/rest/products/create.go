package products

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/products"
	"omg/api/internal/model"

	"github.com/gin-gonic/gin"
)

type createRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       string `json:"price" binding:"required"`
	Stock       string `json:"stock" binding:"required"`
}

type createResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
	Status      string `json:"status"`
}

// Create handles product creates
func (h *Handler) Create(c *gin.Context) {
	var req createRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	price, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stock, err := strconv.ParseInt(req.Stock, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := model.CreateProductInput{
		Name:  req.Name,
		Desc:  req.Description,
		Price: price,
		Stock: stock,
	}

	p, err := h.controller.Create(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, products.ErrProductAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "product already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, createResponse{
		ID:          strconv.FormatInt(p.ID, 10),
		Name:        p.Name,
		Description: p.Description,
		Price:       strconv.FormatFloat(p.Price, 'f', 2, 64),
		Stock:       strconv.FormatInt(stock, 10),
		Status:      p.Status.String(),
	})
}
