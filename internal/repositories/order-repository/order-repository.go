package orderrepository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"hash/fnv"
	"rapidEx/internal/domain/order"
	"rapidEx/internal/domain/stock"

	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Set(ctx context.Context, order *order.Order) error
	All(ctx context.Context, s *stock.Stock) ([]*order.Order, error)
	Del(ctx context.Context, order *order.Order) error
}

type rsClient struct {
	rc *redis.Client
}

func (r *rsClient) Set(ctx context.Context, order *order.Order) error {
	const op = "orderRepository.Set"

	status := r.rc.HSet(ctx, hashTicker(order.Ticker), order.OrderUUID.String(), order)
	if status.Err() != nil {
		return fmt.Errorf("%s: %w", op, status.Err())
	}

	return nil
}

func (r *rsClient) All(ctx context.Context, s *stock.Stock) ([]*order.Order, error) {
	const op = "orderRepository.All"
	stringCmd := r.rc.HGetAll(ctx, hashTicker(s.Ticker))

	switch {
	case stringCmd.Err() == redis.Nil:
		return nil, errors.New("orders not found")
	case stringCmd.Err() != nil:
		return nil, fmt.Errorf("%s: %w", op, stringCmd.Err())
	}

	result, err := stringCmd.Result()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	orders := make([]*order.Order, 0)
	for _, v := range result {
		var o order.Order
		err := o.UnmarshalBinary([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		orders = append(orders, &o)
	}
	return orders, nil
}

func (r *rsClient) Del(ctx context.Context, order *order.Order) error {
	intCmd := r.rc.HDel(ctx, "orders", order.OrderUUID.String())
	log.Println(intCmd.Result())

	if intCmd.Err() != nil {
		return intCmd.Err()
	}

	return nil
}

func NewOrderRepository(rc *redis.Client) Repository {
	return &rsClient{rc: rc}
}


func hashTicker(ticker string) string {
	hash := fnv.New32a()
	hash.Write([]byte(ticker))
	return fmt.Sprintf("%x", hash.Sum(nil))
}