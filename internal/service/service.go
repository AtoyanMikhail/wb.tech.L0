package service

import (
	"context"

	"L0/internal/cache"
	"L0/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *repository.Order) error
	GetOrderByID(ctx context.Context, orderUID string) (*repository.Order, error)
}

type OrderServiceImpl struct {
	repo  repository.OrderRepository
	cache cache.Cache
}

func NewOrderService(repo repository.OrderRepository, cache cache.Cache) OrderService {
	return &OrderServiceImpl{
		repo:  repo,
		cache: cache,
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, order *repository.Order) error {
	// Валидируем заказ
	if err := repository.ValidateOrder(order); err != nil {
		return err
	}

	// Сохраняем в БД
	if err := s.repo.SaveOrder(ctx, order); err != nil {
		return err
	}

	// Сохраняем в кэш
	if err := s.cache.Set(ctx, order.OrderUID, order); err != nil {
		// Логируем ошибку кэша, но не прерываем выполнение
		// TODO: добавить логгер
	}

	return nil
}

func (s *OrderServiceImpl) GetOrderByID(ctx context.Context, orderUID string) (*repository.Order, error) {
	// Сначала проверяем кэш
	if order, err := s.cache.Get(ctx, orderUID); err == nil && order != nil {
		return order, nil
	}

	// Если нет в кэше, получаем из БД
	order, err := s.repo.GetOrderByID(ctx, orderUID)
	if err != nil {
		return nil, err
	}

	if order != nil {
		// Сохраняем в кэш
		if err := s.cache.Set(ctx, orderUID, order); err != nil {
			// Логируем ошибку кэша, но не прерываем выполнение
			// TODO: добавить логгер
		}
	}

	return order, nil
}
