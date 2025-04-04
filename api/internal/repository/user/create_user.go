package user

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository/generator"
	"omg/api/internal/repository/orm"

	"github.com/friendsofgo/errors"
	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// CreateUser saves user in DB
func (i impl) CreateUser(ctx context.Context, m model.User) (model.User, error) {
	id, err := generator.UserIDSNF.Generate()
	if err != nil {
		return model.User{}, pkgerrors.WithStack(err)
	}

	o := orm.User{
		ID:       id,
		Name:     m.Name,
		Email:    m.Email,
		Password: m.Password,
		Status:   m.Status.String(),
	}

	if err = o.Insert(ctx, i.dbConn, boil.Infer()); err != nil {
		return model.User{}, errors.WithStack(err)
	}

	m.ID = id
	m.CreatedAt = o.CreatedAt
	m.UpdatedAt = o.UpdatedAt

	return m, nil
}
