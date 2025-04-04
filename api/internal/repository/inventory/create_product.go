package inventory

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository/generator"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// CreateProduct saves product in DB
func (i impl) CreateProduct(ctx context.Context, p model.Product) (model.Product, error) {
	id, err := generator.ProductIDSNF.Generate()
	if err != nil {
		return p, pkgerrors.WithStack(err)
	}

	o := orm.Product{
		ID:          id,
		Name:        p.Name,
		Description: p.Description,
		Status:      p.Status.String(),
		Price:       p.Price,
		Stock:       p.Stock,
	}

	if err = o.Insert(ctx, i.dbConn, boil.Infer()); err != nil {
		return p, pkgerrors.WithStack(err)
	}

	p.ID = o.ID
	p.CreatedAt = o.CreatedAt
	p.UpdatedAt = o.UpdatedAt

	return p, nil
}
