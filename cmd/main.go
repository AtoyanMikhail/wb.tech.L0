package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"L0/internal/config"
	"L0/internal/repository"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg := config.NewConfig()

	repo, err := repository.NewPostgresRepository(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Получаем абсолютный путь к папке migrations
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatalf("failed to get migrations path: %v", err)
	}

	// Запуск миграций
	err = repo.RunMigrations(migrationsPath)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	fmt.Println("Миграции успешно применены!")

	// Пример заказа
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

	err = repo.SaveOrder(context.Background(), order)
	if err != nil {
		log.Fatalf("failed to save order: %v", err)
	}
	fmt.Println("Тестовый заказ успешно добавлен!")
}

func mustMarshal(v interface{}) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
