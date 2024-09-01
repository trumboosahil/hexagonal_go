package kafka

import (
	"context"
	"fmt"
	"hexagonal_go/internal/domain/entities"
	"hexagonal_go/internal/ports/outbound"
	"log"

	segmentio_kafka "github.com/segmentio/kafka-go"
)

type KafkaOrderRepository struct {
	writer *segmentio_kafka.Writer
}

func NewKafkaOrderRepository(writer *segmentio_kafka.Writer) outbound.OrderRepository {
	return &KafkaOrderRepository{writer: writer}
}

func (r *KafkaOrderRepository) Save(order *entities.Order) error {
	msg := segmentio_kafka.Message{
		Key:   []byte(order.ID),
		Value: []byte(fmt.Sprintf("Amount: %f, Status: %s", order.Amount, order.Status)),
	}
	// Attempt to write the message to Kafka
	if err := r.writer.WriteMessages(context.Background(), msg); err != nil {
		log.Printf("Kafka write error: %v", err)
		return err
	}
	log.Println("Message successfully written to Kafka")
	return nil
}
