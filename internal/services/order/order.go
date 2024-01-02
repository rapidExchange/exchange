package order

import (
	"context"
	"fmt"
	"log/slog"
	"rapidEx/internal/domain/order"
)

type OrdersProvider interface {
	GetAll(ctx context.Context) ([]*order.Order, error)
}

type OrderSaver interface {
	Set(ctx context.Context, order *order.Order) error
}

type OrderMonitor struct {
	log            *slog.Logger
	ordersProvider OrdersProvider
	orderSaver     OrderSaver
}

func New(log *slog.Logger, ordersProvider OrdersProvider, orderSaver OrderSaver) *OrderMonitor {
	return &OrderMonitor{log: log,
		ordersProvider: ordersProvider,
		orderSaver:     orderSaver}
}

func (o *OrderMonitor) GetAll(ctx context.Context) ([]*order.Order, error) {
	const op = "orderMonitor.GetAll"

	log := o.log.With(slog.String("op", op))

	orders, err := o.ordersProvider.GetAll(ctx)
	if err != nil {
		log.Warn("Unable to get orders")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return orders, nil
}
