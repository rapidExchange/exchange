package orderrepository

import (
	"context"
	"errors"
	"log"
	"rapidEx/internal/domain/order"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Set(ctx context.Context, order *order.Order) error
	GetAll(ctx context.Context) ([]*order.Order, error)
	Del(ctx context.Context, order *order.Order) error
}

type rsClient struct {
	rc *redis.Client
}

func (r *rsClient) Set(ctx context.Context, order *order.Order) error {
	status := r.rc.HSet(ctx, "orders", uuid.New().String(), order)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *rsClient) GetAll(ctx context.Context) ([]*order.Order, error) {
	stringCmd := r.rc.HGetAll(ctx, "orders")

	switch {
	case stringCmd.Err() == redis.Nil:
		return nil, errors.New("orders not found")
	case stringCmd.Err() != nil:
		return nil, stringCmd.Err()
	}

	result, err := stringCmd.Result()
	if err != nil {
		return nil, err
	}
	orders := make([]*order.Order, 0)
	for _, v := range result {
		var o order.Order
		err := o.UnmarshalBinary([]byte(v))
		if err != nil {
			return nil, err
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
