package model

import "time"

// OrderStatus represents the status of the order
type OrderStatus string

const (
	// OrderStatusPending means the order is pending
	OrderStatusPending OrderStatus = "PENDING"
	// OrderStatusPaid means the order is paid
	OrderStatusPaid OrderStatus = "PAID"
	// OrderStatusProcessing means the order is processing to delivery
	OrderStatusProcessing OrderStatus = "PROCESSING"
	// OrderStatusShipped means the order is delivery for shipper
	OrderStatusShipped OrderStatus = "SHIPPED"
	// OrderStatusDelivered means the order is delivered
	OrderStatusDelivered OrderStatus = "DELIVERED"
	// OrderStatusCancelled means the order is cancelled
	OrderStatusCancelled OrderStatus = "CANCELLED"
	// OrderStatusFailed means the order is failed when payment
	OrderStatusFailed OrderStatus = "FAILED"
	// OrderStatusRefunded means the order is refunded for client
	OrderStatusRefunded OrderStatus = "REFUNDED"
)

// String converts to string value
func (p OrderStatus) String() string {
	return string(p)
}

// IsValid checks if plan status is valid
func (p OrderStatus) IsValid() bool {
	switch p {
	case OrderStatusPending, OrderStatusPaid, OrderStatusProcessing, OrderStatusShipped, OrderStatusDelivered, OrderStatusCancelled, OrderStatusFailed, OrderStatusRefunded:
		return true
	}
	return false
}

// Order represents the Order
type Order struct {
	ID         int64
	UserID     int64
	Status     OrderStatus
	TotalCost  float64
	OrderItems []OrderItem
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// CreateOrderInput represents the input when create an order
type CreateOrderInput struct {
	UserID int64
	Items  []CreateOrderItemInput
}

type CreateOrderItemInput struct {
	ProductID int64
	Quantity  int64
}
