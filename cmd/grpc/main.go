package main

import (
	"database/sql"
	"hexagonal_go/internal/adapters/driven/kafka"
	"hexagonal_go/internal/adapters/driven/postgres"
	"hexagonal_go/internal/adapters/driving/grpc"
	"hexagonal_go/internal/domain/services"
	pb "hexagonal_go/pkg/grpc/proto"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
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

	grpcOrderHandler := grpc.NewGRPCOrderHandler(orderService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOrderServiceServer(s, grpcOrderHandler)
	log.Printf("gRPC server listening on %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
