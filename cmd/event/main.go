package main

import (
	"context"
	"database/sql"
	"hexagonal_go/internal/adapters/driven/kafka"
	"hexagonal_go/internal/adapters/driven/postgres"
	"hexagonal_go/internal/adapters/driving/event"
	"hexagonal_go/internal/domain/services"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

func main() {
	db, _ := sql.Open("postgres", "user=postgres password=password dbname=mydb sslmode=disable")
	writer := &kafka.Writer{
		Addr:  kafka.TCP("kafka:9092"),
		Topic: "order_topic",
	}

	orderRepo := postgres.NewPostgresOrderRepository(db)
	kafkaProducer := kafka.NewKafkaMessageProducer(writer)
	orderService := services.NewOrderService(orderRepo, kafkaProducer)

	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	eventHandler := event.NewEventOrderHandler(orderService, rdb)

	ctx := context.Background()
	go eventHandler.Listen(ctx, "order_events")

	// Keep the main goroutine alive
	select {}
}
