package repository

import (
	"context"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

// TODO: add nescessary validation for all fields
type Order struct {
	OrderUID          string          `db:"order_uid" validate:"required,alphanum" json:"order_uid"`
	TrackNumber       string          `db:"track_number" validate:"required" json:"track_number"`
	Entry             string          `db:"entry" validate:"required" json:"entry"`
	Delivery          json.RawMessage `db:"delivery" validate:"required" json:"delivery"`
	Payment           json.RawMessage `db:"payment" validate:"required" json:"payment"`
	Items             json.RawMessage `db:"items" validate:"required" json:"items"`
	Locale            string          `db:"locale" validate:"required" json:"locale"`
	InternalSignature string          `db:"internal_signature" json:"internal_signature"`
	CustomerID        string          `db:"customer_id" validate:"required" json:"customer_id"`
	DeliveryService   string          `db:"delivery_service" json:"meest"`
	ShardKey          string          `db:"shardkey" json:"shardkey"`
	SmID              int             `db:"sm_id" json:"sm_id"`
	DateCreated       string          `db:"date_created" validate:"required,datetime=2006-01-02T15:04:05Z" json:"date_created"`
	OofShard          string          `db:"oof_shard" json:"oof_shard"`
}

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *Order) error
	GetOrderByID(ctx context.Context, orderUID string) (*Order, error)
	GetAllOrders(ctx context.Context) ([]Order, error)
}

var validate = validator.New()

func ValidateOrder(order *Order) error {
	return validate.Struct(order)
}
