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

// UpdateProduct updates the product in DB
func (i impl) UpdateProduct(ctx context.Context, p model.Product) (model.Product, error) {
	o, err := orm.FindProduct(ctx, i.dbConn, p.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, ErrProductNotFound
		}
		return model.Product{}, pkgerrors.WithStack(err)
	}

	o.Name = p.Name
	o.Description = p.Description
	o.Price = p.Price
	o.Stock = p.Stock
	o.Status = p.Status.String()
	if _, err = o.Update(ctx, i.dbConn, boil.Whitelist(
		orm.ProductColumns.Name,
		orm.ProductColumns.Description,
		orm.ProductColumns.Price,
		orm.ProductColumns.Stock,
		orm.ProductColumns.Status,
		orm.ProductColumns.UpdatedAt,
	)); err != nil {
		return model.Product{}, pkgerrors.WithStack(err)
	}

	p.UpdatedAt = o.UpdatedAt

	return p, nil
}
