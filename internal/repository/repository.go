package repository

import (
	"L0/internal/models"
	"context"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *models.Order) error
	RunMigrations(migrationsPath string) error
	GetOrderByID(ctx context.Context, orderUID string) (*models.Order, error)
	GetAllOrders(ctx context.Context) ([]models.Order, error)
}
