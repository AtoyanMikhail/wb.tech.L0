package repository

import (
	"context"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type Order struct {
	OrderUID          string          `db:"order_uid" validate:"required,alphanum"`
	TrackNumber       string          `db:"track_number" validate:"required"`
	Entry             string          `db:"entry" validate:"required"`
	Delivery          json.RawMessage `db:"delivery" validate:"required"`
	Payment           json.RawMessage `db:"payment" validate:"required"`
	Items             json.RawMessage `db:"items" validate:"required"`
	Locale            string          `db:"locale" validate:"required"`
	InternalSignature string          `db:"internal_signature"`
	CustomerID        string          `db:"customer_id" validate:"required"`
	DeliveryService   string          `db:"delivery_service" validate:"required"`
	ShardKey          string          `db:"shardkey" validate:"required"`
	SmID              int             `db:"sm_id" validate:"required"`
	DateCreated       string          `db:"date_created" validate:"required,datetime=2006-01-02T15:04:05Z"`
	OofShard          string          `db:"oof_shard" validate:"required"`
}

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, orderUID string) (*Order, error)
}

var validate = validator.New()

func ValidateOrder(order *Order) error {
	return validate.Struct(order)
}
