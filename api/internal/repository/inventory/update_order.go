package inventory

import (
	"context"
	"database/sql"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func (i impl) UpdateOrder(ctx context.Context, m model.Order) (model.Order, error) {
	o, err := orm.FindOrder(ctx, i.dbConn, m.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, ErrOrderNotFound
		}

		return model.Order{}, pkgerrors.WithStack(err)
	}

	o.Status = m.Status.String()
	o.TotalCost = m.TotalCost

	if _, err = o.Update(ctx, i.dbConn, boil.Whitelist(
		orm.OrderColumns.Status,
		orm.OrderColumns.TotalCost,
		orm.OrderColumns.UpdatedAt,
	)); err != nil {
		return model.Order{}, pkgerrors.WithStack(err)
	}

	m.UpdatedAt = o.UpdatedAt

	return m, nil
}
