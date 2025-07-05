package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/kafka"
	"L0/internal/logger"
	"L0/internal/repository"
	"L0/internal/server"
	"L0/internal/service"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log := logger.NewLogger().WithField("component", "main")
	log.Info("Starting L0 service")

	cfg := config.NewConfig()
	log.Info("Configuration loaded")

	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Info("Database connection established")

	cache := cache.NewRedisCache(cfg.Redis, log)
	log.Info("Redis cache initialized")

	orderService := service.NewOrderService(repo, cache, log)
	log.Info("Order service initialized")

	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Errorf("Failed to get migrations path: %v", err)
		log.Fatalf("failed to get migrations path: %v", err)
	}

	err = repo.RunMigrations(migrationsPath)
	if err != nil {
		log.Errorf("Failed to run migrations: %v", err)
		log.Fatalf("failed to run migrations: %v", err)
	}
	log.Info("Database migrations applied successfully")

	consumer := kafka.NewConsumer(cfg, orderService, log)
	log.Info("Kafka consumer initialized")

	// Context for shutting kafka down gracefully
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		log.Info("Starting Kafka consumer")
		if err := consumer.Start(ctx); err != nil {
			log.Errorf("Consumer error: %v", err)
		}
	}()

	handler := server.NewHandler(orderService, log)
	appServer := server.NewServer(handler)
	go func() {
		log.Info("Starting HTTP server on :8081")
		if err := appServer.Run(":8081"); err != nil {
			log.Errorf("Failed to start HTTP server: %v", err)
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Info("Received shutdown signal. Gracefully shutting down...")

	cancel()

	// timeout for graceful shutdown
	time.Sleep(3 * time.Second)
	log.Info("Service shutdown completed")
}
