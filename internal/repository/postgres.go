package repository

import (
	"L0/internal/config"
	"L0/internal/models"
	"context"
	"encoding/json"
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

func NewPostgresRepository(cfg *config.Config) (OrderRepository, error) {
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

// OrderDB is a helper struct for reading orders from db
type OrderDB struct {
	OrderUID          string          `db:"order_uid"`
	TrackNumber       string          `db:"track_number"`
	Entry             string          `db:"entry"`
	Delivery          json.RawMessage `db:"delivery"`
	Payment           json.RawMessage `db:"payment"`
	Items             json.RawMessage `db:"items"`
	Locale            string          `db:"locale"`
	InternalSignature string          `db:"internal_signature"`
	CustomerID        string          `db:"customer_id"`
	DeliveryService   string          `db:"delivery_service"`
	ShardKey          string          `db:"shardkey"`
	SmID              int             `db:"sm_id"`
	DateCreated       string          `db:"date_created"`
	OofShard          string          `db:"oof_shard"`
}

// ToModel converts OrderDB into models.Order
func (o *OrderDB) ToModel() (*models.Order, error) {
	var delivery models.Delivery
	var payment models.Payment
	var items []models.Item

	if err := json.Unmarshal(o.Delivery, &delivery); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(o.Payment, &payment); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(o.Items, &items); err != nil {
		return nil, err
	}

	return &models.Order{
		OrderUID:          o.OrderUID,
		TrackNumber:       o.TrackNumber,
		Entry:             o.Entry,
		Delivery:          delivery,
		Payment:           payment,
		Items:             items,
		Locale:            o.Locale,
		InternalSignature: o.InternalSignature,
		CustomerID:        o.CustomerID,
		DeliveryService:   o.DeliveryService,
		ShardKey:          o.ShardKey,
		SmID:              o.SmID,
		DateCreated:       o.DateCreated,
		OofShard:          o.OofShard,
	}, nil
}

// FromModel converts models.Order into OrderDB
func FromModel(order *models.Order) (*OrderDB, error) {
	deliveryJSON, err := json.Marshal(order.Delivery)
	if err != nil {
		return nil, err
	}
	paymentJSON, err := json.Marshal(order.Payment)
	if err != nil {
		return nil, err
	}
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return nil, err
	}

	return &OrderDB{
		OrderUID:          order.OrderUID,
		TrackNumber:       order.TrackNumber,
		Entry:             order.Entry,
		Delivery:          deliveryJSON,
		Payment:           paymentJSON,
		Items:             itemsJSON,
		Locale:            order.Locale,
		InternalSignature: order.InternalSignature,
		CustomerID:        order.CustomerID,
		DeliveryService:   order.DeliveryService,
		ShardKey:          order.ShardKey,
		SmID:              order.SmID,
		DateCreated:       order.DateCreated,
		OofShard:          order.OofShard,
	}, nil
}

func (r *PostgresRepository) SaveOrder(ctx context.Context, order *models.Order) error {
	orderDB, err := FromModel(order)
	if err != nil {
		return err
	}

	query := `INSERT INTO orders (
		order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES (
		:order_uid, :track_number, :entry, :delivery, :payment, :items, :locale, :internal_signature, :customer_id, :delivery_service, :shardkey, :sm_id, :date_created, :oof_shard
	)`

	_, err = r.db.NamedExecContext(ctx, query, orderDB)
	return err
}

func (r *PostgresRepository) GetOrderByID(ctx context.Context, orderUID string) (*models.Order, error) {
	var orderDB OrderDB
	query := `SELECT * FROM orders WHERE order_uid = $1`

	err := r.db.GetContext(ctx, &orderDB, query, orderUID)
	if err != nil {
		return nil, err
	}

	return orderDB.ToModel()
}

func (r *PostgresRepository) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	var ordersDB []OrderDB
	query := `SELECT * FROM orders`

	err := r.db.SelectContext(ctx, &ordersDB, query)
	if err != nil {
		return nil, err
	}

	orders := make([]models.Order, len(ordersDB))
	for i, orderDB := range ordersDB {
		order, err := orderDB.ToModel()
		if err != nil {
			return nil, err
		}
		orders[i] = *order
	}

	return orders, nil
}
