package inventory

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository/generator"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// CreateOrder saves order in DB
func (i impl) CreateOrder(ctx context.Context, m model.Order) (model.Order, error) {
	id, err := generator.OrderIDSNF.Generate()
	if err != nil {
		return m, pkgerrors.WithStack(err)
	}

	o := orm.Order{
		ID:        id,
		UserID:    m.UserID,
		Status:    m.Status.String(),
		TotalCost: m.TotalCost,
	}

	if err = o.Insert(ctx, i.dbConn, boil.Infer()); err != nil {
		return m, pkgerrors.WithStack(err)
	}

	m.ID = o.ID
	m.CreatedAt = o.CreatedAt
	m.UpdatedAt = o.UpdatedAt

	return m, nil
}
