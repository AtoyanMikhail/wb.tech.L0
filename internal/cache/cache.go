package cache

import (
	"context"
	"sync"

	"L0/internal/repository"
)

type Cache interface {
	Set(ctx context.Context, key string, value *repository.Order) error
	Get(ctx context.Context, key string) (*repository.Order, error)
	Delete(ctx context.Context, key string) error
}

type InMemoryCache struct {
	data map[string]*repository.Order
	mu   sync.RWMutex
}

func NewInMemoryCache() Cache {
	return &InMemoryCache{
		data: make(map[string]*repository.Order),
	}
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value *repository.Order) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
	return nil
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (*repository.Order, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if value, exists := c.data[key]; exists {
		return value, nil
	}
	return nil, nil // Возвращаем nil, если ключ не найден
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}
