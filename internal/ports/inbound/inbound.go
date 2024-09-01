package inbound

import "hexagonal_go/internal/domain/entities"

// OrderHandler defines the methods to handle incoming orders.
type OrderHandler interface {
	HandleOrder(order *entities.Order) error
}
