package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"L0/internal/cache"
	"L0/internal/config"
	"L0/internal/kafka"
	"L0/internal/repository"
	"L0/internal/service"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.NewConfig()

	// Инициализация репозитория
	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Инициализация кэша
	cache := cache.NewInMemoryCache()

	// Инициализация сервиса
	orderService := service.NewOrderService(repo, cache)

	// Получаем абсолютный путь к папке migrations
	migrationsPath, err := filepath.Abs("internal/repository/migrations")
	if err != nil {
		log.Fatalf("failed to get migrations path: %v", err)
	}

	// Запуск миграций
	err = repo.RunMigrations(migrationsPath)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	fmt.Println("Миграции успешно применены!")

	// Инициализация Kafka producer
	producer := kafka.NewProducer(cfg)
	defer producer.Close()

	// Инициализация Kafka consumer с сервисом
	consumer := kafka.NewConsumer(cfg, orderService)

	// Создаём контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем consumer в горутине
	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Consumer error: %v", err)
		}
	}()

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Отправляем тестовый заказ в Kafka
	order := &repository.Order{
		OrderUID:    "b563feb7b2b84b6test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: mustMarshal(map[string]interface{}{
			"name":    "Test Testov",
			"phone":   "+9720000000",
			"zip":     "2639809",
			"city":    "Kiryat Mozkin",
			"address": "Ploshad Mira 15",
			"region":  "Kraiot",
			"email":   "test@gmail.com",
		}),
		Payment: mustMarshal(map[string]interface{}{
			"transaction":   "b563feb7b2b84b6test",
			"request_id":    "",
			"currency":      "USD",
			"provider":      "wbpay",
			"amount":        1817,
			"payment_dt":    1637907727,
			"bank":          "alpha",
			"delivery_cost": 1500,
			"goods_total":   317,
			"custom_fee":    0,
		}),
		Items: mustMarshal([]map[string]interface{}{
			{
				"chrt_id":      9934930,
				"track_number": "WBILMTESTTRACK",
				"price":        453,
				"rid":          "ab4219087a764ae0btest",
				"name":         "Mascaras",
				"sale":         30,
				"size":         "0",
				"total_price":  317,
				"nm_id":        2389212,
				"brand":        "Vivienne Sabo",
				"status":       202,
			},
		}),
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		OofShard:          "1",
	}

	// Отправляем заказ в Kafka (валидация будет в сервисе)
	if err := producer.SendOrder(ctx, order); err != nil {
		log.Fatalf("Failed to send order to Kafka: %v", err)
	}
	fmt.Println("Тестовый заказ отправлен в Kafka!")

	// Ждём сигнала для завершения
	<-sigChan
	fmt.Println("\nПолучен сигнал завершения. Завершаем работу...")

	// Отменяем контекст для остановки consumer
	cancel()

	// Даём время на graceful shutdown
	time.Sleep(2 * time.Second)
	fmt.Println("Приложение завершено.")
}

func mustMarshal(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
