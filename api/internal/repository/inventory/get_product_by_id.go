package inventory

import (
	"context"
	"database/sql"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
)

// GetProductByID retrieve product data by product ID
func (i impl) GetProductByID(ctx context.Context, id int64) (model.Product, error) {
	o, err := orm.Products(
		orm.ProductWhere.ID.EQ(id),
	).One(ctx, i.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Product{}, ErrProductNotFound
		}

		return model.Product{}, pkgerrors.WithStack(err)
	}

	return toProduct(o), nil
}
