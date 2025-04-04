package inventory

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

// ProductsFilter holds filters for getting products list
type ProductsFilter struct {
	Status   []model.ProductStatus
	WithLock bool
}

// ListProducts gets a list of products from DB
func (i impl) ListProducts(ctx context.Context) ([]model.Product, error) {
	slice, err := orm.Products(
		orm.ProductWhere.Status.NEQ(model.ProductStatusDeleted.String()),
		qm.OrderBy(orm.ProductColumns.CreatedAt+" DESC"),
	).All(ctx, i.dbConn)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	var result []model.Product
	for _, o := range slice {
		result = append(result, toProduct(o))
	}

	return result, nil
}
