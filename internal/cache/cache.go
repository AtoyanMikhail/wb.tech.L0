package cache

import (
	"context"
	"encoding/json"
	"time"

	"L0/internal/config"
	"L0/internal/logger"
	"L0/internal/models"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value *models.Order) error
	Get(ctx context.Context, key string) (*models.Order, error)
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	client *redis.Client
	prefix string
	ttl    time.Duration
	logger logger.Logger
}

func NewRedisCache(cfg config.RedisConfig, logger logger.Logger) Cache {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return &RedisCache{
		client: client,
		prefix: cfg.Prefix,
		ttl:    time.Duration(cfg.TTL) * time.Second,
		logger: logger.WithField("component", "redis_cache"),
	}
}

func (c *RedisCache) Set(ctx context.Context, key string, value *models.Order) error {
	b, err := json.Marshal(value)
	if err != nil {
		c.logger.Errorf("Failed to marshal order for cache: %v", err)
		return err
	}
	err = c.client.Set(ctx, c.prefix+key, b, c.ttl).Err()
	if err != nil {
		c.logger.Errorf("Failed to set order in cache: %v", err)
	} else {
		c.logger.Infof("Order cached successfully: %s", key)
	}
	return err
}

func (c *RedisCache) Get(ctx context.Context, key string) (*models.Order, error) {
	val, err := c.client.Get(ctx, c.prefix+key).Result()
	if err == redis.Nil {
		c.logger.Debugf("Order not found in cache: %s", key)
		return nil, nil
	}
	if err != nil {
		c.logger.Errorf("Failed to get order from cache: %v", err)
		return nil, err
	}
	var order models.Order
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		c.logger.Errorf("Failed to unmarshal order from cache: %v", err)
		return nil, err
	}
	c.logger.Infof("Order retrieved from cache: %s", key)
	return &order, nil
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	err := c.client.Del(ctx, c.prefix+key).Err()
	if err != nil {
		c.logger.Errorf("Failed to delete order from cache: %v", err)
	} else {
		c.logger.Infof("Order deleted from cache: %s", key)
	}
	return err
}
