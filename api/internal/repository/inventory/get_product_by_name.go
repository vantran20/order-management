package inventory

import (
	"context"
	"database/sql"
	"errors"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
)

// GetProductByName retrieve product data by product name
func (i impl) GetProductByName(ctx context.Context, name string) (model.Product, error) {
	o, err := orm.Products(
		orm.ProductWhere.Name.EQ(name),
	).One(ctx, i.dbConn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Product{}, ErrProductNotFound
		}

		return model.Product{}, pkgerrors.WithStack(err)
	}

	return toProduct(o), nil
}
