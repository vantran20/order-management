package products

import (
	"context"

	"omg/api/internal/model"
)

// List gets a list of products from DB
func (i impl) List(ctx context.Context) ([]model.Product, error) {
	rs, err := i.repo.Inventory().ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	return rs, nil
}
