package order

import (
	"context"
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
	panic("implement me")
}