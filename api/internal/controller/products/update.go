package products

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/inventory"
)

func (i impl) Update(ctx context.Context, inp model.UpdateProductInput) (model.Product, error) {
	// Check if product with this id already exists
	p, err := i.repo.Inventory().GetProductByID(ctx, inp.ID)
	if err != nil {
		if errors.Is(err, inventory.ErrProductNotFound) {
			return model.Product{}, ErrNotFound
		}
		return model.Product{}, err
	}

	productUpToDate, err := i.repo.Inventory().UpdateProduct(ctx, model.Product{
		ID:          p.ID,
		Name:        inp.Name,
		Description: inp.Desc,
		Price:       inp.Price,
		Stock:       inp.Stock,
		Status:      p.Status,
	})
	if err != nil {
		if errors.Is(err, inventory.ErrProductNotFound) {
			return model.Product{}, ErrNotFound
		}
		return model.Product{}, err
	}

	return productUpToDate, nil
}
