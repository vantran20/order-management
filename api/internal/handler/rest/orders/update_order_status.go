package orders

import (
	"errors"
	"net/http"
	"strconv"

	"omg/api/internal/controller/orders"
	"omg/api/internal/handler/ws"
	"omg/api/internal/model"

	"github.com/gin-gonic/gin"
)

type updateOrderRequest struct {
	Status string `json:"status"`
}

// UpdateProduct handles product updating
func (h *Handler) UpdateProduct(c *gin.Context) {
	var req updateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	orderID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.controller.UpdateOrderStatus(c.Request.Context(), orderID, model.OrderStatus(req.Status))
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrOrderNotFound):
			c.JSON(http.StatusConflict, gin.H{"error": "order not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	// Broadcast order status update via WebSocket
	msg := ws.NewOrderStatusMessage(order.ID, order.UserID, order.Status.String(), order.TotalCost)
	if msgBytes, err := msg.ToJSON(); err == nil {
		h.wsHub.Broadcast <- msgBytes
	}

	c.JSON(http.StatusCreated, "Successfully updated order status")
}
