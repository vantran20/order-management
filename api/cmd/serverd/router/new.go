package router

import (
	"context"

	"omg/api/internal/authenticate"
	"omg/api/internal/controller/orders"
	"omg/api/internal/controller/products"
	"omg/api/internal/controller/system"
	"omg/api/internal/controller/users"
	authenticateRestHandler "omg/api/internal/handler/rest/authenticate"
	orderRestHandler "omg/api/internal/handler/rest/orders"
	productRestHandler "omg/api/internal/handler/rest/products"
	userRestHandler "omg/api/internal/handler/rest/users"
	ws2 "omg/api/internal/ws"

	"github.com/gin-gonic/gin"
)

// New creates a new Router instance
func New(
	ctx context.Context,
	corsOrigins []string,
	isGQLIntrospectionOn bool,
	systemCtrl system.Controller,
	productCtrl products.Controller,
	userCtrl users.Controller,
	orderCtrl orders.Controller,
	authService authenticate.AuthService,
	hub ws2.Hub,
) Router {
	return Router{
		ctx:                     ctx,
		corsOrigins:             corsOrigins,
		isGQLIntrospectionOn:    isGQLIntrospectionOn,
		systemCtrl:              systemCtrl,
		productCtrl:             productCtrl,
		productRestHandler:      productRestHandler.New(productCtrl),
		userCtrl:                userCtrl,
		userRestHandler:         userRestHandler.New(userCtrl),
		orderCtrl:               orderCtrl,
		orderRestHandler:        orderRestHandler.NewHandler(orderCtrl, hub),
		authService:             authService,
		authenticateRestHandler: authenticateRestHandler.New(authService),
		engine:                  gin.Default(),
		hub:                     hub,
		wsHandler:               *ws2.NewWebSocketHandler(hub, authService),
	}
}
