package inventory

import (
	"omg/api/internal/model"
	"omg/api/internal/repository/orm"
)

func toProduct(o *orm.Product) model.Product {
	return model.Product{
		ID:          o.ID,
		Name:        o.Name,
		Description: o.Description,
		Stock:       o.Stock,
		Price:       o.Price,
		Status:      model.ProductStatus(o.Status),
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func toOrder(o *orm.Order) model.Order {
	m := model.Order{
		ID:        o.ID,
		UserID:    o.UserID,
		TotalCost: o.TotalCost,
		Status:    model.OrderStatus(o.Status),
	}

	if o.R != nil && o.R.OrderItems != nil && len(o.R.OrderItems) > 0 {
		for _, item := range o.R.OrderItems {
			m.OrderItems = append(m.OrderItems, toOrderItem(item))
		}

	}

	return m
}

func toOrderItem(o *orm.OrderItem) model.OrderItem {
	return model.OrderItem{
		ID:        o.ID,
		OrderID:   o.OrderID,
		ProductID: o.ProductID,
		Quantity:  o.Quantity,
		Price:     o.Price,
	}
}
