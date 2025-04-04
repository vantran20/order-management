package orders

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository"
)

// Controller represents the specification of this pkg
type Controller interface {
	CreateOrder(context.Context, model.CreateOrderInput) (model.Order, error)
	UpdateOrderStatus(context.Context, int64, model.OrderStatus) (model.Order, error)
}

// New initializes a new Controller instance and returns it
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}

type impl struct {
	repo repository.Registry
}
