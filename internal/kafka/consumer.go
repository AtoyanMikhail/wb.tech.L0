package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"L0/internal/config"
	"L0/internal/logger"
	"L0/internal/repository"
	"L0/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader *kafka.Reader
	svc    service.OrderService
	logger logger.Logger
}

func NewConsumer(cfg *config.Config, svc service.OrderService, logger logger.Logger) *Consumer {
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
		logger: logger.WithField("component", "kafka_consumer"),
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	defer c.reader.Close()
	c.logger.Info("Starting Kafka consumer")

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("Kafka consumer stopped")
			return ctx.Err()
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				c.logger.Errorf("Error reading message: %v", err)
				continue
			}

			c.logger.Infof("Received message from Kafka: topic=%s, partition=%d, offset=%d",
				m.Topic, m.Partition, m.Offset)

			if err := c.processMessage(ctx, m); err != nil {
				c.logger.Errorf("Error processing message: %v", err)
				// Never commit offset on error, so the message will be processed again
				continue
			}

			// Commit offset only after successful processing
			if err := c.reader.CommitMessages(ctx, m); err != nil {
				c.logger.Errorf("Error committing message: %v", err)
			} else {
				c.logger.Infof("Message committed: topic=%s, partition=%d, offset=%d",
					m.Topic, m.Partition, m.Offset)
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, m kafka.Message) error {
	var order repository.Order
	if err := json.Unmarshal(m.Value, &order); err != nil {
		c.logger.Errorf("Failed to unmarshal order: %v", err)
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	c.logger.Infof("Processing order: %s", order.OrderUID)

	// Validation + Save to DB + cache
	if err := c.svc.CreateOrder(ctx, &order); err != nil {
		c.logger.Errorf("Failed to create order: %v", err)
		return fmt.Errorf("failed to create order: %w", err)
	}

	c.logger.Infof("Successfully processed order: %s", order.OrderUID)
	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
