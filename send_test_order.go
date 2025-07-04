package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Читаем тестовый заказ
	orderData, err := ioutil.ReadFile("test_order.json")
	if err != nil {
		log.Fatalf("Failed to read test_order.json: %v", err)
	}

	// Создаем Kafka writer
	w := &kafka.Writer{
		Addr:  kafka.TCP("localhost:9092"),
		Topic: "orders",
	}
	defer w.Close()

	// Создаем сообщение
	msg := kafka.Message{
		Key:   []byte("test-order"),
		Value: orderData,
		Time:  time.Now(),
	}

	// Отправляем сообщение
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = w.WriteMessages(ctx, msg)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Println("Test order sent to Kafka successfully!")
	fmt.Printf("Order UID: %s\n", "b563feb7b2b84b6test")
	fmt.Println("Check your service logs to see if it was processed.")
}
