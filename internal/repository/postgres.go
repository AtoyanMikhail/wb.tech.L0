package repository

import (
	"L0/internal/config"
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(cfg *config.Config) (*PostgresRepository, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.DBName,
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) RunMigrations(migrationsPath string) error {
	driver, err := postgres.WithInstance(r.db, &postgres.Config{})
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
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
	) ON CONFLICT (order_uid) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Delivery,
		order.Payment,
		order.Items,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	)
	return err
}

func (r *PostgresRepository) GetOrderByID(ctx context.Context, orderUID string) (*Order, error) {
	query := `SELECT order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM orders WHERE order_uid = $1`
	row := r.db.QueryRowContext(ctx, query, orderUID)
	var order Order
	err := row.Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&order.Delivery,
		&order.Payment,
		&order.Items,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerID,
		&order.DeliveryService,
		&order.ShardKey,
		&order.SmID,
		&order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}
