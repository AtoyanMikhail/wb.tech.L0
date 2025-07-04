package service

import (
	"context"

	"L0/internal/cache"
	"L0/internal/logger"
	"L0/internal/repository"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order *repository.Order) error
	GetOrderByID(ctx context.Context, orderUID string) (*repository.Order, error)
}

type OrderServiceImpl struct {
	repo   repository.OrderRepository
	cache  cache.Cache
	logger logger.Logger
}

func NewOrderService(repo repository.OrderRepository, cache cache.Cache, logger logger.Logger) OrderService {
	orders, err := repo.GetAllOrders(context.Background())
	if err == nil {
		for _, order := range orders {
			cache.Set(context.Background(), order.OrderUID, &order)
		}
	}

	return &OrderServiceImpl{
		repo:   repo,
		cache:  cache,
		logger: logger.WithField("component", "order_service"),
	}
}

func (s *OrderServiceImpl) CreateOrder(ctx context.Context, order *repository.Order) error {
	s.logger.Infof("Creating order: %s", order.OrderUID)

	if err := repository.ValidateOrder(order); err != nil {
		s.logger.Errorf("Order validation failed: %v", err)
		return err
	}

	if err := s.repo.SaveOrder(ctx, order); err != nil {
		s.logger.Errorf("Failed to save order to database: %v", err)
		return err
	}
	s.logger.Infof("Order saved to database: %s", order.OrderUID)

	if err := s.cache.Set(ctx, order.OrderUID, order); err != nil {
		s.logger.Warnf("Failed to cache order: %v", err)
	}

	return nil
}

func (s *OrderServiceImpl) GetOrderByID(ctx context.Context, orderUID string) (*repository.Order, error) {
	s.logger.Infof("Getting order by ID: %s", orderUID)

	if order, err := s.cache.Get(ctx, orderUID); err == nil && order != nil {
		s.logger.Infof("Order found in cache: %s", orderUID)
		return order, nil
	}

	order, err := s.repo.GetOrderByID(ctx, orderUID)
	if err != nil {
		s.logger.Errorf("Failed to get order from database: %v", err)
		return nil, err
	}

	if order != nil {
		s.logger.Infof("Order found in database: %s", orderUID)
		if err := s.cache.Set(ctx, orderUID, order); err != nil {
			s.logger.Warnf("Failed to cache order: %v", err)
		}
	} else {
		s.logger.Warnf("Order not found: %s", orderUID)
	}

	return order, nil
}
