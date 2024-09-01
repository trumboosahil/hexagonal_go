package grpc

import (
	"context"
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/ports/inbound"
	pb "hexagonal_go/pkg/grpc/proto"
)

type GRPCOrderHandler struct {
	orderService inbound.OrderHandler
}

func NewGRPCOrderHandler(orderService inbound.OrderHandler) *GRPCOrderHandler {
	return &GRPCOrderHandler{orderService: orderService}
}

func (h *GRPCOrderHandler) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	order := &entities.Order{
		ID:     req.Id,
		Amount: req.Amount,
		Status: req.Status,
	}

	err := h.orderService.HandleOrder(order)
	if err != nil {
		return nil, err
	}

	return &pb.OrderResponse{
		Success: true,
	}, nil
}
