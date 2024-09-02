package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"

	_ "github.com/lib/pq"
	segmentio_kafka "github.com/segmentio/kafka-go"

	"hexagonal_go/internal/adapters/driven/kafka"
	"hexagonal_go/internal/adapters/driven/order_redis"
	"hexagonal_go/internal/adapters/driven/postgres"
	"hexagonal_go/internal/adapters/driving/api"
	"hexagonal_go/internal/domain/services"
	"hexagonal_go/internal/ports/outbound"
)

func connectToKafka(broker string) error {
	conn, err := net.Dial("tcp", broker)
	if err != nil {
		// Log the detailed error for debugging
		log.Printf("Error connecting to Kafka broker at %s: %v", broker, err)

		// Return a simple error message with a possible solution
		return fmt.Errorf("failed to connect to Kafka broker at %s. Please ensure the broker is running and accessible on the network", broker)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to Kafka broker at", broker)
	return nil
}

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5434 user=user password=password dbname=mydb sslmode=disable")
	if err != nil {
		panic(err)
	}
	err1 := connectToKafka("localhost:9092")
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}

	writer := segmentio_kafka.Writer{
		Addr:  segmentio_kafka.TCP("localhost:9092"),
		Topic: "order_topic",
	}
	// Initialize Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	var orderRepo outbound.OrderRepository
	useKafka := false
	useRedis := false

	// Choose which repository to use based on your needs
	if useKafka {
		orderRepo = kafka.NewKafkaOrderRepository(&writer)
	} else if useRedis {
		orderRepo = order_redis.NewRedisOrderRepository(rdb, "user_stream")
	} else {
		orderRepo = postgres.NewPostgresOrderRepository(db)
	}

	orderService := services.NewOrderService(orderRepo)

	apiHandler := api.NewAPIOrderHandler(orderService)

	http.HandleFunc("/order", apiHandler.HandleOrder)

	log.Println("API server is listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
