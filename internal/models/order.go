package models

import (
	"github.com/go-playground/validator/v10"
)

type Order struct {
	OrderUID          string   `db:"order_uid" validate:"required,alphanum" json:"order_uid"`
	TrackNumber       string   `db:"track_number" validate:"required" json:"track_number"`
	Entry             string   `db:"entry" validate:"required" json:"entry"`
	Delivery          Delivery `db:"delivery" validate:"required" json:"delivery"`
	Payment           Payment  `db:"payment" validate:"required" json:"payment"`
	Items             []Item   `db:"items" validate:"required" json:"items"`
	Locale            string   `db:"locale" validate:"required" json:"locale"`
	InternalSignature string   `db:"internal_signature" json:"internal_signature"`
	CustomerID        string   `db:"customer_id" validate:"required" json:"customer_id"`
	DeliveryService   string   `db:"delivery_service" json:"delivery_service"`
	ShardKey          string   `db:"shardkey" json:"shardkey"`
	SmID              int      `db:"sm_id" json:"sm_id"`
	DateCreated       string   `db:"date_created" validate:"required,datetime=2006-01-02T15:04:05Z" json:"date_created"`
	OofShard          string   `db:"oof_shard" json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Zip     string `json:"zip" validate:"required"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
	PaymentDt    int64   `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"required"`
	GoodsTotal   float64 `json:"goods_total" validate:"required"`
	CustomFee    float64 `json:"custom_fee"`
}

type Item struct {
	ChrtID      int     `json:"chrt_id" validate:"required"`
	TrackNumber string  `json:"track_number" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
	Rid         string  `json:"rid" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Sale        float64 `json:"sale" validate:"required"`
	Size        string  `json:"size" validate:"required"`
	TotalPrice  float64 `json:"total_price" validate:"required"`
	NmID        int     `json:"nm_id" validate:"required"`
	Brand       string  `json:"brand" validate:"required"`
	Status      int     `json:"status" validate:"required"`
}

var validate = validator.New()

func ValidateOrder(order *Order) error {
	return validate.Struct(order)
}
