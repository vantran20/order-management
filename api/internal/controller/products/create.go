package products

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/inventory"
)

// Create creates the product
func (i impl) Create(ctx context.Context, inp model.CreateProductInput) (model.Product, error) {
	// Check if product with this name already exists
	_, err := i.repo.Inventory().GetProductByName(ctx, inp.Name)
	if err != nil {
		if !errors.Is(err, inventory.ErrProductNotFound) {
			return model.Product{}, err
		}
	} else {
		return model.Product{}, ErrProductAlreadyExists
	}

	product := model.Product{
		Name:        inp.Name,
		Description: inp.Desc,
		Status:      model.ProductStatusActive,
		Price:       inp.Price,
		Stock:       inp.Stock,
	}

	return i.repo.Inventory().CreateProduct(ctx, product)
}
