package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"L0/internal/config"
	"L0/internal/repository"
	"L0/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	svc    service.OrderService
}

func NewConsumer(cfg *config.Config, svc service.OrderService) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    cfg.Kafka.Topic,
		GroupID:  cfg.Kafka.GroupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &Consumer{
		reader: reader,
		svc:    svc,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	defer c.reader.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}

			if err := c.processMessage(ctx, m); err != nil {
				log.Printf("Error processing message: %v", err)
				// Never commit offset on error, so the message will be processed again
				continue
			}

			// Commit offset only after successful processing
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				log.Printf("Error committing message: %v", err)
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, m kafka.Message) error {
	var order repository.Order
	if err := json.Unmarshal(m.Value, &order); err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	// Validation + Save to DB + cache
	if err := c.svc.CreateOrder(ctx, &order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	log.Printf("Successfully processed order: %s", order.OrderUID)
	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
