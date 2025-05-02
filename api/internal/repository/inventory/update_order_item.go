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

// UpdateOrderItem updates the order item in DB
func (i impl) UpdateOrderItem(ctx context.Context, m model.OrderItem) (model.OrderItem, error) {
	o, err := orm.FindOrderItem(ctx, i.dbConn, m.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.OrderItem{}, ErrOrderItemNotFound
		}

		return model.OrderItem{}, pkgerrors.WithStack(err)
	}

	o.Quantity = m.Quantity
	o.Price = m.Price

	if _, err = o.Update(ctx, i.dbConn, boil.Whitelist(
		orm.OrderItemColumns.Quantity,
		orm.OrderItemColumns.Price,
		orm.OrderItemColumns.UpdatedAt,
	)); err != nil {
		return model.OrderItem{}, pkgerrors.WithStack(err)
	}

	m.UpdatedAt = o.UpdatedAt

	return m, nil
}
