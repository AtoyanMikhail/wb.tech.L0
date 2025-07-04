package repository

import (
	"L0/internal/config"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(cfg *config.Config) (*PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) RunMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(r.db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (r *PostgresRepository) SaveOrder(ctx context.Context, order *Order) error {
	query := `INSERT INTO orders (
		order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES (
		:order_uid, :track_number, :entry, :delivery, :payment, :items, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard
	)`

	_, err := r.db.NamedExecContext(ctx, query, order)

	return err
}

func (r *PostgresRepository) GetOrderByID(ctx context.Context, orderUID string) (*Order, error) {
	var order Order
	query := `SELECT * FROM orders WHERE order_uid = $1`

	err := r.db.GetContext(ctx, &order, query, orderUID)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *PostgresRepository) GetAllOrders(ctx context.Context) ([]Order, error) {
	var orders []Order
	query := `SELECT * FROM orders`
	
	err := r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}
	
	return orders, nil
}