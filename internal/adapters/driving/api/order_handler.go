package api

import (
	"encoding/json"
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/domain/services"
	"net/http"
)

type APIOrderHandler struct {
	orderService *services.OrderService
}

func NewAPIOrderHandler(orderService *services.OrderService) *APIOrderHandler {
	return &APIOrderHandler{
		orderService: orderService,
	}
}

func (h *APIOrderHandler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	var order entities.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.orderService.ProcessOrder(&order); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
