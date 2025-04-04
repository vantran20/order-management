package ws

import (
	"encoding/json"
	"strconv"
)

type MessageType string

const (
	MessageTypeOrderStatus MessageType = "order_status"
)

type OrderStatusMessage struct {
	Type      MessageType `json:"type"`
	OrderID   string      `json:"order_id"`
	UserID    string      `json:"user_id"`
	Status    string      `json:"status"`
	TotalCost string      `json:"total_cost"`
}

func NewOrderStatusMessage(orderID, userID int64, status string, totalCost float64) *OrderStatusMessage {
	return &OrderStatusMessage{
		Type:      MessageTypeOrderStatus,
		OrderID:   strconv.FormatInt(orderID, 10),
		UserID:    strconv.FormatInt(userID, 10),
		TotalCost: strconv.FormatFloat(totalCost, 'f', -1, 64),
		Status:    status,
	}
}

func (m *OrderStatusMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
