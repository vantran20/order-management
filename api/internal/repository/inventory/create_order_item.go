package inventory

import (
	"context"

	"omg/api/internal/model"
	"omg/api/internal/repository/generator"
	"omg/api/internal/repository/orm"

	pkgerrors "github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

// CreateOrderItem saves order item in DB
func (i impl) CreateOrderItem(ctx context.Context, m model.OrderItem) (model.OrderItem, error) {
	id, err := generator.OrderItemIDSNF.Generate()
	if err != nil {
		return m, pkgerrors.WithStack(err)
	}

	o := orm.OrderItem{
		ID:        id,
		OrderID:   m.OrderID,
		ProductID: m.ProductID,
		Quantity:  m.Quantity,
		Price:     m.Price,
	}

	if err = o.Insert(ctx, i.dbConn, boil.Infer()); err != nil {
		return m, pkgerrors.WithStack(err)
	}

	m.ID = o.ID
	m.CreatedAt = o.CreatedAt
	m.UpdatedAt = o.UpdatedAt

	return m, nil
}
