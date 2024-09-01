package outbound

import "hexagonal_go/internal/domain/entities"

// OrderRepository defines the methods for interacting with orders.
type OrderRepository interface {
	Save(order *entities.Order) error
}
