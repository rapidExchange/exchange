package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Repository interface {
	Set(ctx context.Context, order Order) error
	Get(ctx context.Context, ticker string) (*Order, error)
	Del(ctx context.Context, ticker string) error
}

type rsClient struct {
	rc *redis.Client
}

func (r *rsClient) Set(ctx context.Context, order Order) error {
	status := r.rc.HSet(ctx, "orders", uuid.New().String(), order)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *rsClient) Get(ctx context.Context, ticker string) (*Order, error) {
	stringCmd := r.rc.HGetAll(ctx, "orders")
	stringCmd.Scan()
	switch {
	case stringCmd.Err() == redis.Nil:
		return nil, errors.New("Order not found")
	case stringCmd.Err() != nil:
		return nil, stringCmd.Err()
	}

	order := Order{}

	err := stringCmd.Scan(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *rsClient) Del(ctx context.Context, ticker string) error {
	intCmd := r.rc.Del(ctx, ticker)

	if intCmd.Err() != nil {
		return intCmd.Err()
	}

	return nil
}

func NewRepository(rc *redis.Client) Repository {
	return &rsClient{rc: rc}
}
