package users

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository"
)

type Controller interface {
	Create(context.Context, model.CreateUserInput) (model.User, error)
	GetByID(context.Context, int64) (model.User, error)
	GetByEmail(context.Context, string) (model.User, error)
	GetUsers(context.Context) ([]model.User, error)
	Update(context.Context, model.UpdateUserInput) (model.User, error)
	Delete(context.Context, int64) error
}

type impl struct {
	repo repository.Registry
}

func New(repo repository.Registry) Controller {
	return &impl{repo: repo}
}
