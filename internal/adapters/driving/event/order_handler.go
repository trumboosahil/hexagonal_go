package event

import (
	"context"
	"encoding/json"
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/ports/inbound"
	"log"

	"github.com/go-redis/redis/v8"
)

type EventOrderHandler struct {
	orderService inbound.OrderHandler
	redisClient  *redis.Client
}

func NewEventOrderHandler(orderService inbound.OrderHandler, redisClient *redis.Client) *EventOrderHandler {
	return &EventOrderHandler{
		orderService: orderService,
		redisClient:  redisClient,
	}
}

func (h *EventOrderHandler) Listen(ctx context.Context, channel string) {
	pubsub := h.redisClient.Subscribe(ctx, channel)
	ch := pubsub.Channel()

	for msg := range ch {
		var order entities.Order
		if err := json.Unmarshal([]byte(msg.Payload), &order); err != nil {
			log.Printf("Failed to unmarshal order: %v", err)
			continue
		}

		err := h.orderService.HandleOrder(&order)
		if err != nil {
			log.Printf("Failed to handle order: %v", err)
		} else {
			log.Printf("Order processed successfully: %v", order.ID)
		}
	}
}
