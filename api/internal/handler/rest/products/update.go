package products

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/products"
	"omg/api/internal/model"
	"omg/api/pkg/floatutil"

	"github.com/gin-gonic/gin"
	pkgerrors "github.com/pkg/errors"
)

type updateProductRequest struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
	Status      string `json:"status"`
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

	input, err := validateAndMapUpdateInput(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.controller.Update(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, products.ErrNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		case errors.Is(err, products.ErrProductDeleted):
			c.JSON(http.StatusBadRequest, gin.H{"error": "product deleted"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, updateProductResponse{
		ID:          strconv.FormatInt(p.ID, 10),
		Name:        p.Name,
		Description: p.Description,
		Price:       floatutil.FormatFloat(p.Price),
		Stock:       strconv.FormatInt(p.Stock, 10),
		Status:      p.Status.String(),
	})
}

func validateAndMapUpdateInput(req updateProductRequest) (model.UpdateProductInput, error) {
	if req.ID == "" {
		return model.UpdateProductInput{}, errors.New("product id is required")
	}

	if req.Name == "" {
		return model.UpdateProductInput{}, errors.New("product name is required")
	}

	if req.Description == "" {
		return model.UpdateProductInput{}, errors.New("product description is required")
	}

	if req.Price == "" {
		return model.UpdateProductInput{}, errors.New("product price is required")
	}

	if req.Stock == "" {
		return model.UpdateProductInput{}, errors.New("product stock is required")
	}

	if req.Status == "" {
		return model.UpdateProductInput{}, errors.New("product status is required")
	}

	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		return model.UpdateProductInput{}, pkgerrors.WithStack(err)
	}

	if id == 0 {
		return model.UpdateProductInput{}, errors.New("invalid product id")
	}

	price, err := strconv.ParseFloat(req.Price, 64)
	if err != nil {
		return model.UpdateProductInput{}, pkgerrors.WithStack(err)
	}

	stock, err := strconv.ParseInt(req.Stock, 10, 64)
	if err != nil {
		return model.UpdateProductInput{}, pkgerrors.WithStack(err)
	}

	status := model.ProductStatus(req.Status)
	if !status.IsValid() {
		return model.UpdateProductInput{}, errors.New("invalid product status")
	}

	return model.UpdateProductInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       price,
		Stock:       stock,
		Status:      model.ProductStatus(req.Status),
	}, nil
}
