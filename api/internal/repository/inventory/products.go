package inventory

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
)

const (
	cacheKeyObjTypeActiveProductsCount = "active-products-count"
	cacheKeyActiveProductsCount        = "all"
)

// GetActiveProductsCountFromDB gets active products count from DB
func (i impl) GetActiveProductsCountFromDB(ctx context.Context) (int64, error) {
	result, err := orm.Products(orm.ProductWhere.Status.EQ(model.ProductStatusActive.String())).Count(ctx, i.dbConn)
	return result, pkgerrors.WithStack(err)
}
