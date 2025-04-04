package orders

import "errors"

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrGetProduct        = errors.New("fail to get product")
	ErrProductOutOfStock = errors.New("product out of stock")
	ErrUpdateProduct     = errors.New("fail to update product")
	ErrCreateOrderItem   = errors.New("fail to create order item")
	ErrCreateOrder       = errors.New("fail to create order")
	ErrUpdateOrder       = errors.New("fail to update order")
	ErrOrderNotFound     = errors.New("order not found")
)
