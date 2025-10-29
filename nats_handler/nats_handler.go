package nats_handler

import (
	"context"
	"log"

	"labs/l0/services"

	"github.com/nats-io/nats.go"
)

type OrderHandler struct {
	orderServce *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderServce: orderService,
	}
}

func (h *OrderHandler) HandlerOrderMessage() nats.MsgHandler {
	return func(msg *nats.Msg) {
		ctx := context.Background()

		log.Printf("Received message from NATS: %s", string(msg.Data))

		if err := h.orderServce.ProcessMessage(ctx, msg.Data); err != nil {
			log.Printf("Error processing the message: %e", err)
			return
		}
	}
}
