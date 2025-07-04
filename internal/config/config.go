package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres PostgresConfig
	Kafka    KafkaConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	Prefix   string
	TTL      int // seconds
}

func NewConfig() *Config {
	godotenv.Load()

	redisDB := 0
	if dbStr := getEnv("REDIS_DB", "0"); dbStr != "" {
		fmt.Sscanf(dbStr, "%d", &redisDB)
	}
	redisTTL := 3600
	if ttlStr := getEnv("REDIS_TTL", "3600"); ttlStr != "" {
		fmt.Sscanf(ttlStr, "%d", &redisTTL)
	}

	return &Config{
		Postgres: PostgresConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DBName:   getEnv("POSTGRES_DB", "orders"),
		},
		Kafka: KafkaConfig{
			Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			Topic:   getEnv("KAFKA_TOPIC", "orders"),
			GroupID: getEnv("KAFKA_GROUP_ID", "orders-service"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
			Prefix:   getEnv("REDIS_PREFIX", "order:"),
			TTL:      redisTTL,
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
