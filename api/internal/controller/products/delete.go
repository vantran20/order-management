package products

import (
	"context"
	"errors"
	"time"

	"omg/api/internal/model"
	"omg/api/internal/repository"
	"omg/api/internal/repository/inventory"
	"omg/api/pkg/db/pg"
)

// Delete deletes a product from DB by marking the status as deleted
func (i impl) Delete(ctx context.Context, id int64) error {
	txFunc := func(newCtx context.Context, repo repository.Registry) error {
		p, err := repo.Inventory().GetProductByID(newCtx, id)
		if err != nil {
			if errors.Is(err, inventory.ErrProductNotFound) {
				return ErrNotFound
			}
			return err
		}

		if p.Status == model.ProductStatusDeleted {
			return ErrProductDeleted
		}

		p.Status = model.ProductStatusDeleted
		if _, err = repo.Inventory().UpdateProduct(newCtx, p); err != nil {
			if errors.Is(err, inventory.ErrProductNotFound) {
				return ErrNotFound
			}
			return err
		}

		return nil
	}

	// Create a new context with timeout for the transaction
	newCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Use the new context with timeout for the transaction
	return i.repo.DoInTx(newCtx, txFunc, pg.ExponentialBackOff(2, 2*time.Minute))
}
