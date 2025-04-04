package router

import (
	"context"
	"net/http"

	"omg/api/internal/authenticate"
	"omg/api/internal/controller/orders"
	"omg/api/internal/controller/products"
	"omg/api/internal/controller/system"
	"omg/api/internal/controller/users"
	authenticateRestHandler "omg/api/internal/handler/rest/authenticate"
	orderRestHandler "omg/api/internal/handler/rest/orders"
	productRestHandler "omg/api/internal/handler/rest/products"
	userRestHandler "omg/api/internal/handler/rest/users"
	"omg/api/internal/handler/ws"

	"github.com/gin-gonic/gin"
)

type RouteGroup func(*gin.RouterGroup)

// Router defines the routes & handlers of the app
type Router struct {
	ctx                     context.Context
	corsOrigins             []string
	isGQLIntrospectionOn    bool
	systemCtrl              system.Controller
	productCtrl             products.Controller
	productRestHandler      productRestHandler.Handler
	userCtrl                users.Controller
	userRestHandler         userRestHandler.Handler
	orderCtrl               orders.Controller
	orderRestHandler        orderRestHandler.Handler
	authService             authenticate.AuthService
	authenticateRestHandler authenticateRestHandler.Handler
	engine                  *gin.Engine
	wsHandler               ws.WebSocketHandler
	hub                     *ws.Hub
}

// Handler returns the Handler for use by the server
func (rtr *Router) Handler() http.Handler {
	// Start the WebSocket hub
	go rtr.hub.Run()

	// Set up CORS
	rtr.engine.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	rtr.setupRoutes(rtr.engine)

	return rtr.engine
}

func (rtr *Router) setupRoutes(r *gin.Engine) {
	public := r.Group("/public")
	rtr.public(public)

	authenticated := r.Group("/authenticated")
	authenticated.Use(authenticate.NewAuthMiddleware(rtr.authService).Handler())
	rtr.authenticated(authenticated)
}

func (rtr *Router) public(rg *gin.RouterGroup) {
	usersRouter := rg.Group("/users")
	usersRouter.POST("/register", rtr.userRestHandler.Register)
	usersRouter.POST("/login", rtr.authenticateRestHandler.Login)
	usersRouter.GET("/ws", rtr.wsHandler.Handle)
}

func (rtr *Router) authenticated(rg *gin.RouterGroup) {
	usersRouter := rg.Group("/users")
	usersRouter.GET("/profile", rtr.userRestHandler.GetUserByEmail)
	usersRouter.GET("/:id", rtr.userRestHandler.GetUserByID)
	usersRouter.GET("/list", rtr.userRestHandler.List)
	usersRouter.POST("/update", rtr.userRestHandler.UpdateUser)

	productsRouter := rg.Group("/products")
	productsRouter.POST("/create", rtr.productRestHandler.Create)
	productsRouter.POST("/update", rtr.productRestHandler.UpdateProduct)
	productsRouter.POST("/delete", rtr.productRestHandler.Delete)
	productsRouter.GET("/:id", rtr.productRestHandler.GetProductByID)
	productsRouter.GET("/list", rtr.productRestHandler.List)

	orderRouter := rg.Group("/order")
	orderRouter.POST("/create", rtr.orderRestHandler.Create)
	orderRouter.POST("/update/:id", rtr.orderRestHandler.UpdateProduct)
	orderRouter.GET("/ws", rtr.wsHandler.HandleOrderUpdates)
}
