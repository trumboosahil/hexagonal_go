package postgres

import (
	"database/sql"
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/ports/outbound"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) outbound.OrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Save(order *entities.Order) error {
	_, err := r.db.Exec("INSERT INTO orders (order_id, amount, status) VALUES ($1, $2, $3)",
		order.ID, order.Amount, order.Status)
	return err
}
