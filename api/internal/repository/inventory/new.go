package inventory

import (
	"context"

	"omg/api/internal/model"
	"omg/api/pkg/db/pg"
)

// Repository provides the specification of the functionality provided by this pkg
type Repository interface {
	ListProducts(context.Context) ([]model.Product, error)
	CreateProduct(context.Context, model.Product) (model.Product, error)
	UpdateProduct(context.Context, model.Product) (model.Product, error)
	GetProductByName(context.Context, string) (model.Product, error)
	GetProductByID(context.Context, int64) (model.Product, error)

	CreateOrder(context.Context, model.Order) (model.Order, error)
	CreateOrderItem(context.Context, model.OrderItem) (model.OrderItem, error)
	UpdateOrder(context.Context, model.Order) (model.Order, error)
	UpdateOrderItem(context.Context, model.OrderItem) (model.OrderItem, error)
	GetOrderByID(context.Context, int64) (model.Order, error)
}

// New returns an implementation instance satisfying Repository
func New(dbConn pg.ContextExecutor) Repository {
	return impl{dbConn: dbConn}
}

type impl struct {
	dbConn pg.ContextExecutor
}
