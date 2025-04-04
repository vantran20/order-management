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

type createOrderRequest struct {
	UserID string `json:"user_id"`
	Items  []struct {
		ProductID string `json:"product_id"`
		Quantity  string `json:"quantity"`
	} `json:"items"`
}

type createOrderResponse struct {
	ID        string                    `json:"id"`
	UserID    string                    `json:"user_id"`
	TotalCost string                    `json:"total_cost"`
	Status    string                    `json:"status"`
	Items     []createOrderItemResponse `json:"items"`
}

type createOrderItemResponse struct {
	ID        string `json:"id"`
	OrderId   string `json:"order_id"`
	ProductID string `json:"product_id"`
	Quantity  string `json:"quantity"`
	Price     string `json:"price"`
}

// Create handles product creates
func (h *Handler) Create(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := strconv.ParseInt(req.UserID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if userID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	input := model.CreateOrderInput{
		UserID: userID,
	}

	for _, item := range req.Items {
		productID, err := strconv.ParseInt(item.ProductID, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if productID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
			return
		}
		quantity, err := strconv.ParseInt(item.Quantity, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if quantity == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "quantity required"})
			return
		}
		input.Items = append(input.Items, model.CreateOrderItemInput{
			ProductID: productID,
			Quantity:  quantity,
		})
	}

	order, err := h.controller.CreateOrder(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, orders.ErrProductNotFound):
			c.JSON(http.StatusConflict, gin.H{"error": "product not found"})
		case errors.Is(err, orders.ErrGetProduct):
			c.JSON(http.StatusConflict, gin.H{"error": "fail to get product"})
		case errors.Is(err, orders.ErrCreateOrder):
			c.JSON(http.StatusConflict, gin.H{"error": "fail to create order"})
		case errors.Is(err, orders.ErrCreateOrderItem):
			c.JSON(http.StatusConflict, gin.H{"error": "fail to create order item"})
		case errors.Is(err, orders.ErrProductOutOfStock):
			c.JSON(http.StatusConflict, gin.H{"error": "product out of stock"})
		case errors.Is(err, orders.ErrUpdateProduct):
			c.JSON(http.StatusConflict, gin.H{"error": "fail to update product"})
		case errors.Is(err, orders.ErrCreateOrder):
			c.JSON(http.StatusConflict, gin.H{"error": "fail to create order"})
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

	resp := createOrderResponse{
		ID:        strconv.FormatInt(order.ID, 10),
		UserID:    strconv.FormatInt(order.UserID, 10),
		TotalCost: strconv.FormatFloat(order.TotalCost, 'f', -1, 64),
		Status:    order.Status.String(),
	}

	for _, item := range order.OrderItems {
		resp.Items = append(resp.Items, createOrderItemResponse{
			ID:        strconv.FormatInt(item.ID, 10),
			OrderId:   strconv.FormatInt(item.OrderID, 10),
			ProductID: strconv.FormatInt(item.ProductID, 10),
			Quantity:  strconv.FormatInt(item.Quantity, 10),
			Price:     strconv.FormatFloat(item.Price, 'f', -1, 64),
		})
	}
	c.JSON(http.StatusCreated, resp)
}
