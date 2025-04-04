package inventory

import (
	"context"
	"database/sql"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// GetOrderByID retrieve order data by order ID
func (i impl) GetOrderByID(ctx context.Context, id int64) (model.Order, error) {
	o, err := orm.Orders(
		orm.OrderWhere.ID.EQ(id),
		qm.Load(orm.OrderRels.OrderItems),
	).One(ctx, i.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Order{}, ErrOrderNotFound
		}

		return model.Order{}, pkgerrors.WithStack(err)
	}

	return toOrder(o), nil
}
