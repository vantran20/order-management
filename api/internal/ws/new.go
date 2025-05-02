package ws

import (
	"sync"

	"omg/api/internal/authenticate"

	"github.com/gin-gonic/gin"
)

type WebSocket interface {
	Handle(c *gin.Context)
	HandleOrderUpdates(c *gin.Context)
}

func NewWebSocketHandler(hub Hub, authService authenticate.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
	}
}

type WebSocketHandler struct {
	hub         Hub
	authService authenticate.AuthService
}

type Hub interface {
	Run()
	Register(client *Client)
	Unregister(client *Client)
	BroadcastMessage(message []byte)
}

func NewHub() Hub {
	return &implHub{
		Broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

type implHub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}
