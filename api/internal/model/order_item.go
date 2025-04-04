package model

import "time"

// OrderItem represents the order item
type OrderItem struct {
	ID        int64
	OrderID   int64
	ProductID int64
	Quantity  int64
	Price     float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
