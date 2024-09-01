package main

import (
	"database/sql"
	"hexagonal_go/internal/adapters/driven/kafka"
	"hexagonal_go/internal/adapters/driven/postgres"
	"hexagonal_go/internal/adapters/driving/websocket"
	"hexagonal_go/internal/domain/services"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

	wsHandler := websocket.NewWebSocketOrderHandler(orderService)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			log.Printf("Failed to upgrade connection: %v", err)
			return
		}
		go wsHandler.HandleOrder(conn)
	})

	log.Println("WebSocket server started on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalf("Failed to start WebSocket server: %v", err)
	}
}
