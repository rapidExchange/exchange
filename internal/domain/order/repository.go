package order

import (
	"context"
	"errors"

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
	status := r.rc.Set(ctx, order.Ticker, order, 1)
	if status.Err() != nil {
		return status.Err()
	}

	return nil
}

func (r *rsClient) Get(ctx context.Context, ticker string) (*Order, error) {
	order := r.rc.Get(ctx, ticker)
	switch {
	case order.Err() == redis.Nil:
		return nil, errors.New("Order not found")
	case order.Err() != nil:
		return nil, order.Err()
	}

	s := Order{}

	err := order.Scan(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
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
