package websocket

import (
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/ports/inbound"
	"log"

	"github.com/gorilla/websocket"
)

type WebSocketOrderHandler struct {
	orderService inbound.OrderHandler
}

func NewWebSocketOrderHandler(orderService inbound.OrderHandler) *WebSocketOrderHandler {
	return &WebSocketOrderHandler{orderService: orderService}
}

func (h *WebSocketOrderHandler) HandleOrder(conn *websocket.Conn) {
	defer conn.Close()

	for {
		var order entities.Order
		if err := conn.ReadJSON(&order); err != nil {
			log.Printf("Error reading JSON: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid order data"))
			return
		}

		err := h.orderService.HandleOrder(&order)
		if err != nil {
			log.Printf("Error handling order: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to process order"))
			return
		}

		conn.WriteMessage(websocket.TextMessage, []byte("Order processed successfully"))
	}
}
