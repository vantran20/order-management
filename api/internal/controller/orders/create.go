package orders

import (
	"context"
	"errors"
	"time"

	"omg/api/internal/model"
	"omg/api/internal/repository"
	"omg/api/internal/repository/inventory"
	"omg/api/pkg/db/pg"
)

func (i impl) CreateOrder(ctx context.Context, inp model.CreateOrderInput) (model.Order, error) {
	var order model.Order

	txFunc := func(newCtx context.Context, repo repository.Registry) error {
		var err error
		order, err = i.processOrder(newCtx, repo, inp)
		return err
	}

	// Create a new context with timeout for the transaction
	newCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Use the new context with timeout for the transaction
	if err := i.repo.DoInTx(newCtx, txFunc, pg.ExponentialBackOff(2, 2*time.Minute)); err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (i impl) processOrder(ctx context.Context, repo repository.Registry, inp model.CreateOrderInput) (model.Order, error) {
	// Create order first
	order := model.Order{
		UserID:    inp.UserID,
		Status:    model.OrderStatusPending,
		TotalCost: 0, // Will be updated after processing items
	}

	order, err := repo.Inventory().CreateOrder(ctx, order)
	if err != nil {
		return model.Order{}, ErrCreateOrder
	}

	// Process items with the created order ID
	items, totalCost, err := i.processOrderItems(ctx, repo, order.ID, inp.Items)
	if err != nil {
		return model.Order{}, err
	}

	// Update order with total cost
	order.TotalCost = totalCost
	order, err = repo.Inventory().UpdateOrder(ctx, order)
	if err != nil {
		return model.Order{}, ErrUpdateOrder
	}

	order.OrderItems = items
	return order, nil
}

func (i impl) processOrderItems(ctx context.Context, repo repository.Registry, orderID int64, items []model.CreateOrderItemInput) ([]model.OrderItem, float64, error) {
	var totalCost float64
	var processedItems []model.OrderItem

	for _, item := range items {
		orderItem, itemCost, err := i.processOrderItem(ctx, repo, orderID, item)
		if err != nil {
			return nil, 0, err
		}

		processedItems = append(processedItems, orderItem)
		totalCost += itemCost
	}

	return processedItems, totalCost, nil
}

func (i impl) processOrderItem(ctx context.Context, repo repository.Registry, orderID int64, item model.CreateOrderItemInput) (model.OrderItem, float64, error) {
	// Find product
	product, err := repo.Inventory().GetProductByID(ctx, item.ProductID)
	if err != nil {
		if errors.Is(err, inventory.ErrProductNotFound) {
			return model.OrderItem{}, 0, ErrProductNotFound
		}
		return model.OrderItem{}, 0, ErrGetProduct
	}

	// Check enough stock
	if product.Stock < item.Quantity {
		return model.OrderItem{}, 0, ErrProductOutOfStock
	}

	// Deduct stock of product
	product.Stock -= item.Quantity

	// Update stock for product
	if _, err = repo.Inventory().UpdateProduct(ctx, product); err != nil {
		if errors.Is(err, inventory.ErrProductNotFound) {
			return model.OrderItem{}, 0, ErrProductNotFound
		}
		return model.OrderItem{}, 0, ErrUpdateProduct
	}

	orderItem := model.OrderItem{
		OrderID:   orderID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     product.Price,
	}

	// Insert order item
	if _, err = repo.Inventory().CreateOrderItem(ctx, orderItem); err != nil {
		return model.OrderItem{}, 0, ErrCreateOrderItem
	}

	itemCost := float64(item.Quantity) * product.Price
	return orderItem, itemCost, nil
}
