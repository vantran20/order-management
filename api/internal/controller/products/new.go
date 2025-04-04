package products

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository"
)

// Controller represents the specification of this pkg
type Controller interface {
	List(context.Context) ([]model.Product, error)
	GetByID(context.Context, int64) (model.Product, error)
	Create(context.Context, model.CreateProductInput) (model.Product, error)
	Delete(context.Context, int64) error
	Update(context.Context, model.UpdateProductInput) (model.Product, error)
}

// New initializes a new Controller instance and returns it
func New(repo repository.Registry) Controller {
	return impl{repo: repo}
}

type impl struct {
	repo repository.Registry
}
