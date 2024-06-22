package order

import (
	"context"
	"fmt"
	"log/slog"
	"rapidEx/internal/domain/order"
	"rapidEx/internal/domain/stock"
)

type OrdersProvider interface {
	All(ctx context.Context, s *stock.Stock) ([]*order.Order, error)
}

type OrderSaver interface {
	Set(ctx context.Context, order *order.Order) error
}

type OrderDeleter interface {
	Del(ctx context.Context, order *order.Order) error
}

type OrderService struct {
	log            *slog.Logger
	ordersProvider OrdersProvider
	orderSaver     OrderSaver
	orderDeleter   OrderDeleter
}

func New(log *slog.Logger, ordersProvider OrdersProvider, orderSaver OrderSaver, orderDeleter OrderDeleter) *OrderService {
	return &OrderService{log: log,
		ordersProvider: ordersProvider,
		orderSaver:     orderSaver,
		orderDeleter:   orderDeleter}
}

func (o *OrderService) Del(ctx context.Context, order *order.Order) error {
	const op = "orderMonitor.Del"

	log := o.log.With(slog.String("op", op), slog.String("uuid", fmt.Sprintf("%s\t %s\t %s\t %f", order.Email,
		order.Ticker, order.OrderUUID, order.Quantity)))

	err := o.orderDeleter.Del(ctx, order)
	if err != nil {
		log.Warn("Unable to delete order")
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (o *OrderService) All(ctx context.Context, s *stock.Stock) ([]*order.Order, error) {
	const op = "orderMonitor.GetAll"

	log := o.log.With(slog.String("op", op))

	orders, err := o.ordersProvider.All(ctx, s)
	if err != nil {
		log.Warn("Unable to get orders")
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return orders, nil
}
