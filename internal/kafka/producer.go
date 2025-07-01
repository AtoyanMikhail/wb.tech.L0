package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"L0/internal/config"
	"L0/internal/repository"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg *config.Config) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Brokers...),
		Topic:        cfg.Kafka.Topic,
		RequiredAcks: kafka.RequireAll, // Подтверждение от всех реплик
		Async:        false,            // Синхронная отправка
		Logger:       kafka.LoggerFunc(log.Printf),
	}

	return &Producer{
		writer: writer,
	}
}

func (p *Producer) SendOrder(ctx context.Context, order *repository.Order) error {
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("failed to marshal order: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(order.OrderUID),
		Value: orderBytes,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
