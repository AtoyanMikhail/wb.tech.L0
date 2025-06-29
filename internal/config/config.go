package config

import "os"

type Config struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
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
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
