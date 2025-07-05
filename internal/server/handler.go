package server

import (
	"encoding/json"
	"net/http"

	"L0/internal/logger"
	"L0/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	orderService service.OrderService
	logger       logger.Logger
}

func NewHandler(orderService service.OrderService, logger logger.Logger) *Handler {
	return &Handler{
		orderService: orderService,
		logger:       logger.WithField("component", "http_handler"),
	}
}

func (h *Handler) GetOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		orderUID := c.Param("order_uid")
		h.logger.Infof("HTTP request: GET /order/%s", orderUID)

		if orderUID == "" {
			h.logger.Warn("Empty order_uid in request")
			c.JSON(http.StatusBadRequest, gin.H{"error": "order_uid is required"})
			return
		}

		order, err := h.orderService.GetOrderByID(c.Request.Context(), orderUID)
		if err != nil {
			h.logger.Errorf("Failed to get order %s: %v", orderUID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get order by id"})
			return
		}
		if order == nil {
			h.logger.Warnf("Order not found: %s", orderUID)
			c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
			return
		}

		var delivery, payment interface{}
		var items []interface{}
		if err := json.Unmarshal(order.Delivery, &delivery); err != nil {
			h.logger.Warnf("Failed to unmarshal delivery for order %s: %v", orderUID, err)
			delivery = nil
		}
		if err := json.Unmarshal(order.Payment, &payment); err != nil {
			h.logger.Warnf("Failed to unmarshal payment for order %s: %v", orderUID, err)
			payment = nil
		}
		if err := json.Unmarshal(order.Items, &items); err != nil {
			h.logger.Warnf("Failed to unmarshal items for order %s: %v", orderUID, err)
			items = nil
		}

		response := gin.H{
			"order_uid":          order.OrderUID,
			"track_number":       order.TrackNumber,
			"entry":              order.Entry,
			"delivery":           delivery,
			"payment":            payment,
			"items":              items,
			"locale":             order.Locale,
			"internal_signature": order.InternalSignature,
			"customer_id":        order.CustomerID,
			"delivery_service":   order.DeliveryService,
			"shardkey":           order.ShardKey,
			"sm_id":              order.SmID,
			"date_created":       order.DateCreated,
			"oof_shard":          order.OofShard,
		}

		h.logger.Infof("Order %s returned successfully", orderUID)
		c.JSON(http.StatusOK, response)
	}
}
