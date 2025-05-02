package ws

import (
	"encoding/json"
	"log"
	"strconv"
)

func (h *implHub) Run() {
	log.Printf("Starting WebSocket hub")

	for {
		select {
		case client := <-h.register:
			log.Printf("New client registered: %v", client)
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			log.Printf("Client unregistered: %v", client)
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered. Total clients: %d", len(h.clients))

		case message := <-h.Broadcast:
			log.Printf("Broadcasting message: %+v", message)
			h.mu.RLock()
			for client := range h.clients {
				log.Printf("Starting...")
				// If message is an order status update, check user ID
				if msg, err := parseOrderStatusMessage(message); err == nil {
					log.Printf("Broadcasting message: %+v", msg)
					// Send to all clients if userID is 0, or only to matching user
					userID, err := strconv.ParseInt(msg.UserID, 10, 64)
					if err == nil {
						log.Printf("Failed to parse data %v", msg.UserID)
						close(client.send)
						delete(h.clients, client)
					}

					if client.userID == 0 || client.userID == userID {
						log.Printf("Sending order status update to user %d - OrderID: %s, Status: %s",
							client.userID, msg.OrderID, msg.Status)
						select {
						case client.send <- message:
						default:
							log.Printf("Failed to send message to client %v", client)
							close(client.send)
							delete(h.clients, client)
						}
					}
				} else {
					// For non-order messages, send to all clients
					select {
					case client.send <- message:
					default:
						log.Printf("Failed to send message to client %v", client)
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func parseOrderStatusMessage(data []byte) (*OrderStatusMessage, error) {
	var msg OrderStatusMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func (h *implHub) BroadcastMessage(message []byte) {
	h.Broadcast <- message
}

func (h *implHub) Register(client *Client) {
	h.register <- client
}

func (h *implHub) Unregister(client *Client) {
	h.unregister <- client
}
