package services

import (
	"errors"
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/ports/outbound"
)

type OrderService struct {
	orderRepo outbound.OrderRepository
}

func NewOrderService(orderRepo outbound.OrderRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo}
}

func (s *OrderService) ProcessOrder(order *entities.Order) error {
	if order.Amount <= 0 {
		return errors.New("invalid order amount")
	}

	return s.orderRepo.Save(order)
}
