package orders

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/orders"
	"omg/api/internal/model"
	"omg/api/internal/ws"

	"github.com/gin-gonic/gin"
)

type updateOrderRequest struct {
	Status string `json:"status"`
}

type updateOrderResponse struct {
	ID        string                    `json:"id"`
	UserID    string                    `json:"user_id"`
	TotalCost string                    `json:"total_cost"`
	Status    string                    `json:"status"`
	Items     []updateOrderItemResponse `json:"items"`
}

type updateOrderItemResponse struct {
	ID        string `json:"id"`
	OrderId   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  string `json:"quantity"`
	Price     string `json:"price"`
}

// UpdateOrderStatus handles order updating status
func (h *Handler) UpdateOrderStatus(c *gin.Context) {

	var req updateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if orderID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	order, err := h.controller.UpdateOrderStatus(c.Request.Context(), orderID, model.OrderStatus(req.Status))
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrOrderNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "order not found"})
		case errors.Is(err, orders.ErrProductNotFound):
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		case errors.Is(err, orders.ErrProductOutOfStock):
			c.JSON(http.StatusBadRequest, gin.H{"error": "product out of stock"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Broadcast order status update via WebSocket
	msg := ws.NewOrderStatusMessage(order.ID, order.UserID, order.Status.String(), order.TotalCost)
	if msgBytes, err := msg.ToJSON(); err == nil {
		h.wsHub.BroadcastMessage(msgBytes)
	}

	resp := updateOrderResponse{
		ID:        strconv.FormatInt(order.ID, 10),
		UserID:    strconv.FormatInt(order.UserID, 10),
		TotalCost: strconv.FormatFloat(order.TotalCost, 'f', -1, 64),
		Status:    order.Status.String(),
	}

	for _, item := range order.OrderItems {
		resp.Items = append(resp.Items, updateOrderItemResponse{
			ID:        strconv.FormatInt(item.ID, 10),
			OrderId:   strconv.FormatInt(item.OrderID, 10),
			ProductID: strconv.FormatInt(item.ProductID, 10),
			Quantity:  strconv.FormatInt(item.Quantity, 10),
			Price:     strconv.FormatFloat(item.Price, 'f', -1, 64),
		})
	}
	c.JSON(http.StatusCreated, resp)
}
