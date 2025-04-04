package user

import (
	"context"

	"omg/api/internal/model"
	"omg/api/pkg/db/pg"
)

type Repository interface {
	CreateUser(context.Context, model.User) (model.User, error)
	GetByEmail(context.Context, string) (model.User, error)
	GetByID(context.Context, int64) (model.User, error)
	GetUsers(context.Context) ([]model.User, error)
	Update(context.Context, model.User) error
}

type impl struct {
	dbConn pg.ContextExecutor
}

func New(dbConn pg.ContextExecutor) Repository {
	return &impl{dbConn: dbConn}
}
