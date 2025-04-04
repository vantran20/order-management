package ws

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"omg/api/internal/authenticate"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	Subprotocols: []string{""},
}

type WebSocketHandler struct {
	hub         *Hub
	authService authenticate.AuthService
}

func NewWebSocketHandler(hub *Hub, authService authenticate.AuthService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		authService: authService,
	}
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	if err := h.handleWebSocket(c, 0); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *WebSocketHandler) HandleOrderUpdates(c *gin.Context) {
	// Get and validate token
	tokenString, err := h.getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Validate token and get user ID
	claims, err := h.authService.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if err := h.handleWebSocket(c, claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (h *WebSocketHandler) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header required")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return "", errors.New("invalid authorization header format")
	}

	return tokenString, nil
}

func (h *WebSocketHandler) handleWebSocket(c *gin.Context, userID int64) error {
	// Set required headers
	h.setWebSocketHeaders(c)

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return err
	}

	// Create and register client
	client := NewClientWithUserID(h.hub, conn, userID)
	h.hub.register <- client

	// Start message pumps
	go client.writePump()
	go client.readPump()

	return nil
}

func (h *WebSocketHandler) setWebSocketHeaders(c *gin.Context) {
	headers := map[string]string{
		"Connection":            "Upgrade",
		"Upgrade":               "websocket",
		"Sec-Websocket-Version": "13",
		"Sec-Websocket-Key":     "dGhlIHNhbXBsZSBub25jZQ==",
	}

	for key, value := range headers {
		if c.Request.Header.Get(key) == "" {
			c.Request.Header.Set(key, value)
		}
	}
}
