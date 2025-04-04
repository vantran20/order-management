package products

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/products"
	"omg/api/internal/model"

	"github.com/gin-gonic/gin"
)

type updateProductRequest struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
	Price  string `json:"price"`
	Stock  string `json:"stock"`
	Status string `json:"status"`
}

type updateProductResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
	Status      string `json:"status"`
}

// UpdateProduct handles product updating
func (h *Handler) UpdateProduct(c *gin.Context) {
	var req updateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
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

	input := model.UpdateProductInput{
		ID:     id,
		Name:   req.Name,
		Desc:   req.Desc,
		Price:  price,
		Stock:  stock,
		Status: model.ProductStatus(req.Status),
	}

	p, err := h.controller.Update(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusConflict, gin.H{"error": "product not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, updateProductResponse{
		ID:          strconv.FormatInt(p.ID, 10),
		Name:        p.Name,
		Description: p.Description,
		Price:       strconv.FormatFloat(p.Price, 'f', 2, 64),
		Stock:       strconv.FormatInt(stock, 10),
		Status:      p.Status.String(),
	})
}
