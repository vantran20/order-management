package products

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/inventory"
)

// GetByID gets a single product by product id from DB
func (i impl) GetByID(ctx context.Context, id int64) (model.Product, error) {
	p, err := i.repo.Inventory().GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, inventory.ErrProductNotFound) {
			return model.Product{}, ErrNotFound
		}
		return model.Product{}, err
	}

	return p, nil
}
