package repository

import (
	"context"
	"encoding/json"
)

type Order struct {
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

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, orderUID string) (*Order, error)
}
