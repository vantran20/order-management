package orders

import (
	"context"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/inventory"
)

func (i impl) UpdateOrderStatus(ctx context.Context, id int64, status model.OrderStatus) (model.Order, error) {
	// Check if product with this name already exists
	o, err := i.repo.Inventory().GetOrderByID(ctx, id)
	if err != nil {
		if errors.Is(err, inventory.ErrOrderNotFound) {
			return model.Order{}, ErrOrderNotFound
		}

		return model.Order{}, err
	}

	o.Status = status

	// Update status
	rs, err := i.repo.Inventory().UpdateOrder(ctx, o)
	if err != nil {
		if errors.Is(err, inventory.ErrOrderNotFound) {
			return model.Order{}, ErrOrderNotFound
		}
		return model.Order{}, err
	}

	return rs, nil
}
