package config

import "os"

type Config struct {
	Postgres PostgresConfig
	Kafka    KafkaConfig
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

func NewConfig() *Config {
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
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
