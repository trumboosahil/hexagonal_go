package order_redis

import (
	"context"
	"encoding/json"
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/ports/outbound"

	"github.com/go-redis/redis/v8"
)

type RedisOrderRepository struct {
	client *redis.Client
	stream string
}

func NewRedisOrderRepository(client *redis.Client, stream string) outbound.OrderRepository {
	return &RedisOrderRepository{
		client: client,
		stream: stream,
	}
}

func (r *RedisOrderRepository) Save(user *entities.Order) error {
	// Serialize the user data to JSON
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Add the user data to the Redis stream
	err = r.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: r.stream,
		Values: map[string]interface{}{
			"user": string(data),
		},
	}).Err()

	return err
}
