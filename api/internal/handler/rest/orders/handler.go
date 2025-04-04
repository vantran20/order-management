package orders

import (
	"omg/api/internal/controller/orders"
	"omg/api/internal/handler/ws"
)

type Handler struct {
	controller orders.Controller
	wsHub      *ws.Hub
}

func NewHandler(controller orders.Controller, wsHub *ws.Hub) Handler {
	return Handler{
		controller: controller,
		wsHub:      wsHub,
	}
}
